package vst2

import (
	"fmt"
	"time"

	"pipelined.dev/signal"
)

// Processor represents vst2 sound processor
type Processor struct {
	VST
	plugin *Plugin

	bufferSize  int
	numChannels int
	sampleRate  signal.SampleRate

	currentPosition int64

	// references are needed to free them in Flush.
	doubleIn  DoubleBuffer
	doubleOut DoubleBuffer
}

// Process returns processor function with default settings initialized.
func (p *Processor) Process(pipeID string, sampleRate signal.SampleRate, numChannels int) (func(signal.Float64) error, error) {
	p.sampleRate = sampleRate
	p.numChannels = numChannels
	p.plugin = p.VST.Load(p.callback())

	p.plugin.SetSampleRate(int(p.sampleRate))
	p.plugin.SetSpeakerArrangement(newSpeakerArrangement(p.numChannels), newSpeakerArrangement(p.numChannels))

	for _, d := range p.DispatchBeforeStart {
		p.plugin.Dispatch(d.Opcode, d.Index, d.Value, d.Ptr, d.Opt)
	}
	
	p.plugin.Start()
	var size int
	var out signal.Float64
	return func(in signal.Float64) error {
		// new buffer size.
		if size != in.Size() {
			size = in.Size()
			p.plugin.SetBufferSize(size)

			// reset buffers.
			p.doubleIn.Free()
			p.doubleOut.Free()
			p.doubleIn = NewDoubleBuffer(numChannels, size)
			p.doubleOut = NewDoubleBuffer(numChannels, size)
			out = signal.Float64Buffer(numChannels, size)
		}
		p.doubleIn.CopyFrom(in)
		p.plugin.ProcessDouble(p.doubleIn, p.doubleOut)
		p.currentPosition += int64(in.Size())
		p.doubleOut.CopyTo(out)

		// copy result back to input buffer.
		for i := range out {
			copy(in[i], out[i])
		}
		return nil
	}, nil
}

// Flush suspends plugin.
func (p *Processor) Flush(string) error {
	p.plugin.Stop()
	p.doubleIn.Free()
	p.doubleOut.Free()
	return nil
}

// wraped callback with session.
func (p *Processor) callback() HostCallbackFunc {
	return func(opcode HostOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
		switch opcode {
		case HostIdle:
			p.plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)
		case HostGetCurrentProcessLevel:
			return Return(ProcessLevelRealtime)
		case HostGetSampleRate:
			return Return(p.sampleRate)
		case HostGetBlockSize:
			return Return(p.bufferSize)
		case HostGetTime:
			nanoseconds := time.Now().UnixNano()
			ti := &TimeInfo{
				SampleRate:         float64(p.sampleRate),
				SamplePos:          float64(p.currentPosition),
				NanoSeconds:        float64(nanoseconds),
				TimeSigNumerator:   4,
				TimeSigDenominator: 4,
			}
			return ti.Return()
		default:
			// log.Printf("Plugin requested value of opcode %v\n", opcode)
			break
		}
		return 0
	}
}
