package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include <stdlib.h>
#include <stdint.h>
#include "vst.h"
*/
import "C"
import (
	"fmt"
	"path/filepath"
	"sync"
	"unsafe"
)

// global state for plugins.
var (
	mutex   sync.RWMutex
	plugins = make(map[*effect]*Plugin)
)

//export hostCallback
// global hostCallback, calls real callback.
func hostCallback(e *effect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) Return {
	// AudioMasterVersion is requested when plugin is created
	// It's never in map
	if HostOpcode(opcode) == HostVersion {
		return version
	}
	mutex.RLock()
	p, ok := plugins[e]
	mutex.RUnlock()
	if !ok {
		panic("plugin was closed")
	}

	if p == nil || p.callback == nil {
		panic("host callback is undefined")
	}
	return p.callback(p, HostOpcode(opcode), Index(index), Value(value), Ptr(ptr), Opt(opt))
}

const (
	// VST main function name.
	main = "VSTPluginMain"
	// VST API version.
	version = 2400
)

type (
	// entryPoint is a reference to VST main function. It also keeps
	// reference to VST handle to clean up on Close.
	entryPoint struct {
		main effectMain
		// handle is OS-specific.
		handle
	}

	// wrapper on C entry point.
	effectMain C.EntryPoint

	// Index is index in plugin dispatch/host callback.
	Index int64
	// Value is value in plugin dispatch/host callback.
	Value int64
	// Ptr is ptr in plugin dispatch/host callback.
	Ptr unsafe.Pointer
	// Opt is opt in plugin dispatch/host callback.
	Opt float64
	// Return is returned value for dispatch/host callback.
	Return int64
)

type (
	// Effect is an alias on C effect type.
	effect C.Effect

	// HostCallbackFunc used as callback function called by plugin.
	HostCallbackFunc func(*Plugin, HostOpcode, Index, Value, Ptr, Opt) Return

	// VST used to create new instances of plugin.
	VST struct {
		entryPoint entryPoint
		Name       string
		Path       string
	}

	// Plugin is VST2 plugin instance.
	Plugin struct {
		*effect
		Name     string
		Path     string
		callback HostCallbackFunc
	}
)

// Close cleans up VST handle.
func (e entryPoint) Close() error {
	if e.main == nil {
		return nil
	}
	e.main = nil
	return e.handle.close()
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (e *effect) Dispatch(opcode EffectOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
	return Return(C.dispatch((*C.Effect)(e), C.int(opcode), C.int(index), C.int64_t(value), unsafe.Pointer(ptr), C.float(opt)))
}

// CanProcessFloat32 checks if plugin can process float32.
func (e *effect) CanProcessFloat32() bool {
	if e == nil {
		return false
	}
	return EffectFlags(e.flags)&EffFlagsCanReplacing == EffFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64.
func (e *effect) CanProcessFloat64() bool {
	if e == nil {
		return false
	}
	return EffectFlags(e.flags)&EffFlagsCanDoubleReplacing == EffFlagsCanDoubleReplacing
}

// ProcessDouble audio with VST plugin.
func (e *effect) ProcessDouble(in, out DoubleBuffer) {
	C.processDouble(
		(*C.Effect)(e),
		C.int(in.numChannels),
		C.int(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// ProcessFloat32 audio with VST plugin.
// TODO: add c buffer parameter.
func (e *effect) ProcessFloat32(in [][]float32) (out [][]float32) {
	numChannels := len(in)
	blocksize := len(in[0])

	// convert [][]float32 to []*C.float
	input := make([]*C.float, numChannels)
	output := make([]*C.float, numChannels)
	for i, row := range in {
		// allocate input memory for C layout
		inp := (*C.float)(C.malloc(C.size_t(C.sizeof_float * blocksize)))
		input[i] = inp
		defer C.free(unsafe.Pointer(inp))

		// copy data from slice to C array
		pa := (*[1 << 30]C.float)(unsafe.Pointer(inp))
		for j, v := range row {
			(*pa)[j] = C.float(v)
		}

		// allocate output memory for C layout
		outp := (*C.float)(C.malloc(C.size_t(C.sizeof_float * blocksize)))
		output[i] = outp
		defer C.free(unsafe.Pointer(outp))
	}

	C.processFloat((*C.Effect)(e), C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.float slices to [][]float32
	out = make([][]float32, numChannels)
	for i, data := range output {
		// copy data from C array to slice
		pa := (*[1 << 30]C.float)(unsafe.Pointer(data))
		out[i] = make([]float32, blocksize)
		for j := range out[i] {
			out[i][j] = float32(pa[j])
		}
	}
	return out
}

// Start the plugin.
func (e *effect) Start() {
	e.Dispatch(EffStateChanged, 0, 1, nil, 0.0)
}

// Stop the plugin.
func (e *effect) Stop() {
	e.Dispatch(EffStateChanged, 0, 0, nil, 0.0)
}

// SetBufferSize sets a buffer size
func (e *effect) SetBufferSize(bufferSize int) {
	e.Dispatch(EffSetBufferSize, 0, Value(bufferSize), nil, 0.0)
}

// SetSampleRate sets a sample rate for plugin
func (e *effect) SetSampleRate(sampleRate int) {
	e.Dispatch(EffSetSampleRate, 0, 0, nil, Opt(sampleRate))
}

// SetSpeakerArrangement craetes and passes SpeakerArrangement structures to plugin
func (e *effect) SetSpeakerArrangement(in, out *SpeakerArrangement) {
	e.Dispatch(EffSetSpeakerArrangement, 0, in.Value(), out.Ptr(), 0.0)
}

// ScanPaths returns a slice of default vst2 locations.
// Locations are OS-specific.
func ScanPaths() (paths []string) {
	return append([]string{}, scanPaths...)
}

// Open loads the VST into memory and stores entry point func.
func Open(path string) (*VST, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	ep, err := open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to load VST '%s': %w", path, err)
	}

	return &VST{
		Path:       p,
		entryPoint: ep,
	}, nil
}

// Close cleans up VST resoures.
func (v *VST) Close() error {
	if err := v.entryPoint.Close(); err != nil {
		return fmt.Errorf("failed close VST %s: %w", v.Name, err)
	}
	return nil
}

// Load new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcode.
func (v *VST) Load(c HostCallbackFunc) *Plugin {
	e := (*effect)(C.loadEffect(v.entryPoint.main))
	p := &Plugin{
		effect:   e,
		Path:     v.Path,
		Name:     v.Name,
		callback: c,
	}
	mutex.Lock()
	plugins[e] = p
	mutex.Unlock()
	e.Dispatch(EffOpen, 0, 0, nil, 0.0)
	return p
}

// Close cleans up C refs for plugin
func (p *Plugin) Close() error {
	p.Dispatch(EffClose, 0, 0, nil, 0.0)
	p.effect = nil
	return nil
}

func newSpeakerArrangement(numChannels int) *SpeakerArrangement {
	sa := SpeakerArrangement{}
	sa.NumChannels = int32(numChannels)
	switch numChannels {
	case 0:
		sa.Type = SpeakerArrEmpty
	case 1:
		sa.Type = SpeakerArrMono
	case 2:
		sa.Type = SpeakerArrStereo
	case 3:
		sa.Type = SpeakerArr30Music
	case 4:
		sa.Type = SpeakerArr40Music
	case 5:
		sa.Type = SpeakerArr50
	case 6:
		sa.Type = SpeakerArr60Music
	case 7:
		sa.Type = SpeakerArr70Music
	case 8:
		sa.Type = SpeakerArr80Music
	default:
		sa.Type = SpeakerArrUserDefined
	}

	for i := 0; i < int(numChannels); i++ {
		sa.Speakers[i].Type = SpeakerUndefined
	}
	return &sa
}
