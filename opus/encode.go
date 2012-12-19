package opus

// #cgo pkg-config: opus
// #include <opus.h>
import "C"

type (
	Application int
	Error struct {
		Code int
		Message string
	}
	Encoder struct {
		ptr *C.OpusEncoder
	}
)

func (this *Error) Error() string {
	return this.Message
}

const (
	VOIP Application = C.OPUS_APPLICATION_VOIP
	AUDIO Application = C.OPUS_APPLICATION_AUDIO
	RESTRICTED_LOWDELAY Application = C.OPUS_APPLICATION_RESTRICTED_LOWDELAY

	OK = C.OPUS_OK
	BAD_ARG = C.OPUS_BAD_ARG
	BUFFER_TOO_SMALL
)

func toError(code C.int) error {
	err := &Error{ int(code), "Unknown Error" }
	switch code {
	case C.OPUS_OK:
		return nil
	case C.OPUS_BAD_ARG:
		err.Message = "One or more invalid/out of range arguments"
	case C.OPUS_BUFFER_TOO_SMALL:
		err.Message = "The mode struct passed is invalid"
	case C.OPUS_INTERNAL_ERROR:
		err.Message = "An internal error was detected"
	case C.OPUS_INVALID_PACKET:
		err.Message = "The compressed data passed is corrupted"
	case C.OPUS_UNIMPLEMENTED:
		err.Message = "Invalid/unsupported request number"
	case C.OPUS_INVALID_STATE:
		err.Message = "An encoder or decoder structure is invalid or already freed"
	case C.OPUS_ALLOC_FAIL:
		err.Message = "Memory allocation has failed"
	}
	return err
}

func NewEncoder(sampleRate int, channels int, application Application) (*Encoder, error) {
	var err C.int
	encoder := C.opus_encoder_create(
		C.opus_int32(C.int(sampleRate)),
		C.int(channels),
		C.int(application),
		&err,
	)

	return &Encoder{encoder}, toError(err)
}

func (this *Encoder) Init(sampleRate int, channels int, application Application) {
	C.opus_encoder_init(
		this.ptr,
		C.opus_int32(C.int(sampleRate)),
		C.int(channels),
		C.int(application),
	)
}

func (this *Encoder) Close() {
	C.opus_encoder_destroy(this.ptr)
}
