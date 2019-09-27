package vst2

const (
	maxProgNameLen   = 24 // used for #effGetProgramName, #effSetProgramName, #effGetProgramNameIndexed
	maxParamStrLen   = 8  // used for #effGetParamLabel, #effGetParamDisplay, #effGetParamName
	maxVendorStrLen  = 64 // used for #effGetVendorString, #audioMasterGetVendorString
	maxProductStrLen = 64 // used for #effGetProductString, #audioMasterGetProductString
	maxEffectNameLen = 32 // used for #effGetEffectName

	maxNameLen       = 64  // used for #MidiProgramName, #MidiProgramCategory, #MidiKeyName, #VstSpeakerProperties, #VstPinProperties
	maxLabelLen      = 64  // used for #VstParameterProperties->label, #VstPinProperties->label
	maxShortLabelLen = 8   // used for #VstParameterProperties->shortLabel, #VstPinProperties->shortLabel
	maxCategLabelLen = 24  // used for #VstParameterProperties->label
	maxFileNameLen   = 100 // used for #VstAudioFile->name
)

const (
	// EffectMagic is constant in every plugin.
	EffectMagic = "VstP"
)

// TimeInfo describes the time at the start of the block currently being processed.
type (
	TimeInfo struct {
		// Current Position in audio samples.
		SamplePos float64
		// Current Sample Rate in Herz.
		SampleRate float64
		// System Time in nanoseconds.
		NanoSeconds float64
		// Musical Position, in Quarter Note (1.0 equals 1 Quarter Note).
		PpqPos float64
		// Current Tempo in BPM (Beats Per Minute).
		Tempo float64
		// Last Bar Start Position, in Quarter Note.
		BarStartPos float64
		// Cycle Start (left locator), in Quarter Note.
		CycleStartPos float64
		// Cycle End (right locator), in Quarter Note.
		CycleEndPos float64
		// Time Signature Numerator (e.g. 3 for 3/4).
		TimeSigNumerator int32
		// Time Signature Denominator (e.g. 4 for 3/4).
		TimeSigDenominator int32
		// SMPTE offset in SMPTE subframes (bits; 1/80 of a frame).
		// The current SMPTE position can be calculated using SamplePos, SampleRate, and SMPTEFrameRate.
		SMPTEOffset int32
		// SMPTEFrameRate value.
		SMPTEFrameRate
		// MIDI Clock Resolution (24 Per Quarter Note), can be negative (nearest clock).
		SamplesToNextClock int32
		// TimeInfoFlags values.
		Flags TimeInfoFlags
	}

	// TimeInfoFlags used in TimeInfo.
	TimeInfoFlags int32

	// SMPTEFrameRate values, used in TimeInfo.
	SMPTEFrameRate int32
)

const (
	// TransportChanged is set if play, cycle or record state has changed.
	TransportChanged TimeInfoFlags = 1 << iota
	// TransportPlaying is set if Host sequencer is currently playing
	TransportPlaying
	// TransportCycleActive is set if Host sequencer is in cycle mode.
	TransportCycleActive
	// TransportRecording is set if Host sequencer is in record mode.
	TransportRecording
	_
	_
	// AutomationWriting is set if automation write mode active.
	AutomationWriting
	// AutomationReading is set if automation read mode active.
	AutomationReading
	// NanosValid is set if TimeInfo.NanoSeconds are valid.
	NanosValid
	// PpqPosValid is set if TimeInfo.PpqPos is valid.
	PpqPosValid
	// TempoValid is set if TimeInfo.Tempo is valid.
	TempoValid
	// BarsValid is set if TimeInfo.BarStartPos is valid.
	BarsValid
	// CyclePosValid is set if both TimeInfo.CycleStartPos and TimeInfo.CycleEndPos are valid.
	CyclePosValid
	// TimeSigValid is set if both TimeInfo.TimeSigNumerator and TimeInfo.TimeSigDenominator are valid.
	TimeSigValid
	// SMPTEValid is set if both TimeInfo.SMPTEOffset and TimeInfo.SMPTEFrameRate are valid.
	SMPTEValid
	// ClockValid is set if TimeInfo.SamplesToNextClock are valid.
	ClockValid
)

const (
	// SMPTE24fps is 24 fps.
	SMPTE24fps SMPTEFrameRate = iota
	// SMPTE25fps is 25 fps.
	SMPTE25fps
	// SMPTE2997fps is 29.97 fps.
	SMPTE2997fps
	// SMPTE30fps is 30 fps.
	SMPTE30fps
	// SMPTE2997dfps is 29.97 drop.
	SMPTE2997dfps
	// SMPTE30dfps is 30 drop.
	SMPTE30dfps
	// SMPTEFilm16mm is Film 16mm.
	SMPTEFilm16mm
	// SMPTEFilm35mm is Film 35mm.
	SMPTEFilm35mm
	_
	_
	// SMPTE239fps is HDTV 23.976 fps.
	SMPTE239fps
	// SMPTE249fps is HDTV 24.976 fps.
	SMPTE249fps
	// SMPTE599fps is HDTV 59.94 fps.
	SMPTE599fps
	// SMPTE60fps is HDTV 60 fps.
	SMPTE60fps
)

type (
	// SpeakerArrangement contains information about a channel.
	SpeakerArrangement struct {
		Type        SpeakerArrangementType
		NumChannels int32
		Speakers    [8]Speaker
	}

	// SpeakerArrangementType indicates how the channels are intended to be used in the plugin.
	// Only useful for some hosts.
	SpeakerArrangementType int32

	// Speaker configuration.
	Speaker struct {
		Azimuth   float32
		Elevation float32
		Radius    float32
		Reserved  float32
		Name      [64]byte
		Type      SpeakerType
		Future    [28]byte
	}

	// SpeakerType of particular speaker.
	SpeakerType int32
)

const (
	// SpeakerArrUserDefined is user defined.
	SpeakerArrUserDefined SpeakerArrangementType = iota - 2
	// SpeakerArrEmpty is empty arrangement.
	SpeakerArrEmpty
	// SpeakerArrMono is M.
	SpeakerArrMono
	// SpeakerArrStereo is L R.
	SpeakerArrStereo
	// SpeakerArrStereoSurround is Ls Rs.
	SpeakerArrStereoSurround
	// SpeakerArrStereoCenter is Lc Rc.
	SpeakerArrStereoCenter
	// SpeakerArrStereoSide is Sl Sr.
	SpeakerArrStereoSide
	// SpeakerArrStereoCLfe is C Lfe.
	SpeakerArrStereoCLfe
	// SpeakerArr30Cine is L R C.
	SpeakerArr30Cine
	// SpeakerArr30Music is L R S.
	SpeakerArr30Music
	// SpeakerArr31Cine is L R C Lfe.
	SpeakerArr31Cine
	// SpeakerArr31Music is L R Lfe S.
	SpeakerArr31Music
	// SpeakerArr40Cine is L R C S (LCRS).
	SpeakerArr40Cine
	// SpeakerArr40Music is L R Ls Rs (Quadro).
	SpeakerArr40Music
	// SpeakerArr41Cine is L R C Lfe S (LCRS+Lfe).
	SpeakerArr41Cine
	// SpeakerArr41Music is L R Lfe Ls Rs (Quadro+Lfe).
	SpeakerArr41Music
	// SpeakerArr50 is L R C Ls Rs.
	SpeakerArr50
	// SpeakerArr51 is L R C Lfe Ls Rs.
	SpeakerArr51
	// SpeakerArr60Cine is L R C Ls Rs Cs.
	SpeakerArr60Cine
	// SpeakerArr60Music is L R Ls Rs Sl Sr.
	SpeakerArr60Music
	// SpeakerArr61Cine is L R C Lfe Ls Rs Cs.
	SpeakerArr61Cine
	// SpeakerArr61Music is L R Lfe Ls Rs Sl Sr.
	SpeakerArr61Music
	// SpeakerArr70Cine is L R C Ls Rs Lc Rc.
	SpeakerArr70Cine
	// SpeakerArr70Music is L R C Ls Rs Sl Sr.
	SpeakerArr70Music
	// SpeakerArr71Cine is L R C Lfe Ls Rs Lc Rc.
	SpeakerArr71Cine
	// SpeakerArr71Music is L R C Lfe Ls Rs Sl Sr.
	SpeakerArr71Music
	// SpeakerArr80Cine is L R C Ls Rs Lc Rc Cs.
	SpeakerArr80Cine
	// SpeakerArr80Music is L R C Ls Rs Cs Sl Sr.
	SpeakerArr80Music
	// SpeakerArr81Cine is L R C Lfe Ls Rs Lc Rc Cs.
	SpeakerArr81Cine
	// SpeakerArr81Music is L R C Lfe Ls Rs Cs Sl Sr.
	SpeakerArr81Music
	// SpeakerArr102 is L R C Lfe Ls Rs Tfl Tfc Tfr Trl Trr Lfe2.
	SpeakerArr102
	// not defined.
	numSpeakerArr
)

const (
	// SpeakerUndefined is undefined.
	SpeakerUndefined SpeakerType = 0x7fffffff
	// SpeakerM is Mono (M).
	SpeakerM = iota
	// SpeakerL is Left (L).
	SpeakerL
	// SpeakerR is Right (R).
	SpeakerR
	// SpeakerC is Center (C).
	SpeakerC
	// SpeakerLfe is Subbass (Lfe).
	SpeakerLfe
	// SpeakerLs is Left Surround (Ls).
	SpeakerLs
	// SpeakerRs is Right Surround (Rs).
	SpeakerRs
	// SpeakerLc is Left of Center (Lc).
	SpeakerLc
	// SpeakerRc is Right of Center (Rc).
	SpeakerRc
	// SpeakerS is Surround (S).
	SpeakerS
	// SpeakerCs is Center of Surround (Cs) = Surround (S).
	SpeakerCs = SpeakerS
	// SpeakerSl is Side Left (Sl).
	SpeakerSl
	// SpeakerSr is Side Right (Sr).
	SpeakerSr
	// SpeakerTm is Top Middle (Tm).
	SpeakerTm
	// SpeakerTfl is Top Front Left (Tfl).
	SpeakerTfl
	// SpeakerTfc is Top Front Center (Tfc).
	SpeakerTfc
	// SpeakerTfr is Top Front Right (Tfr).
	SpeakerTfr
	// SpeakerTrl is Top Rear Left (Trl).
	SpeakerTrl
	// SpeakerTrc is Top Rear Center (Trc).
	SpeakerTrc
	// SpeakerTrr is Top Rear Right (Trr).
	SpeakerTrr
	// SpeakerLfe2 is Subbass 2 (Lfe2).
	SpeakerLfe2
)

// EffectFlags values.
type EffectFlags int32

const (
	// EffFlagsHasEditor is set if the plugin provides a custom editor.
	EffFlagsHasEditor EffectFlags = 1 << iota
	_
	_
	_
	// EffFlagsCanReplacing is set if plugin supports replacing process mode.
	EffFlagsCanReplacing
	// EffFlagsProgramChunks is set if preset data is handled in formatless chunks.
	EffFlagsProgramChunks
	_
	_
	// EffFlagsIsSynth is set if plugin is a synth.
	EffFlagsIsSynth
	// EffFlagsNoSoundInStop is set if plugin does not produce sound when input is silence.
	EffFlagsNoSoundInStop
	_
	_
	// EffFlagsCanDoubleReplacing is set if plugin supports double precision processing.
	EffFlagsCanDoubleReplacing

	// deprecated in VST v2.4
	effFlagsHasClip
	// deprecated in VST v2.4
	effFlagsHasVu
	// deprecated in VST v2.4
	effFlagsCanMono
	// deprecated in VST v2.4
	effFlagsExtIsAsync
	// deprecated in VST v2.4
	effFlagsExtHasBuffer
)

// ProcessLevels are used as result for in HostGetCurrentProcessLevel call.
// It tells the plugin in which thread host is right now.
type ProcessLevels int32

const (
	// ProcessLevelUnknown is returned when not supported by host.
	ProcessLevelUnknown ProcessLevels = iota
	// ProcessLevelUser is returned when in user thread (GUI).
	ProcessLevelUser
	// ProcessLevelRealtime is returned when in audio thread (where process is called).
	ProcessLevelRealtime
	// ProcessLevelPrefetch is returned when in sequencer thread (MIDI, timer etc).
	ProcessLevelPrefetch
	// ProcessLevelOffline is returned when in offline processing and thus in user thread.
	ProcessLevelOffline
)
