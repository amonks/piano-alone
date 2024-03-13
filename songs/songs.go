package songs

import (
	"bytes"
	_ "embed"

	"gitlab.com/gomidi/midi/v2/smf"
)

var (
	//go:embed excerpt.mid
	ExcerptBytes []byte
	ExcerptSMF   = mustRead(ExcerptBytes)

	//go:embed prelude_opus_3_no_2.mid
	PreludeOpus3No2Bytes []byte
	PreludeOpus3No2SMF   = mustRead(PreludeOpus3No2Bytes)

	//go:embed prelude_bergamasque.mid
	PreludeBergamasqueBytes []byte
	PreludeBergamasqueSMF   = mustRead(PreludeBergamasqueBytes)
)

func mustRead(bs []byte) *smf.SMF {
	r := bytes.NewReader(bs)
	f, err := smf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	return f
}
