package friends

import (
	"context"
	"fmt"
	"strings"

	da "github.com/kechako/go-yahoo-da"
	"golang.org/x/xerrors"
)

// FirstPersons is a slice of first person words.
var FirstPersons = []string{
	"私",
	"わたし",
	"ワタシ",
	"わたくし",
	"ワタクシ",
	"自分",
	"じぶん",
	"ジブン",
	"僕",
	"ぼく",
	"ボク",
	"俺",
	"おれ",
	"オレ",
	"儂",
	"わし",
	"ワシ",
	"あたし",
	"あたし",
	"あたくし",
	"アタクシ",
	"あてくし",
	"アテクシ",
	"あたい",
	"アタイ",
	"わい",
	"ワイ",
	"わて",
	"ワテ",
	"あて",
	"アテ",
	"わだす",
	"ワダス",
	"あだす",
	"アダス",
	"わす",
	"ワス",
	"内",
	"うち",
	"ウチ",
	"己等",
	"おいら",
	"オイラ",
	"俺ら",
	"おら",
	"オラ",
	"おい",
	"オイ",
	"おいどん",
	"オイドン",
	"うら",
	"ウラ",
	"わ",
	"ワ",
	"わー",
	"ワー",
	"ぼくちゃん",
	"ボクチャン",
	"ぼくちん",
	"ボクチン",
	"おれっち",
	"オレッチ",
}

// A Friends represents a generator of Friends' statement.
type Friends struct {
	client *da.Client
}

// New returns a new *Friends.
func New(appID string) *Friends {
	return &Friends{
		client: da.NewClient(appID),
	}
}

// Say returns a Friends' statement that generated from text.
// Say returns a empty string If text has a no specialty.
func (f *Friends) Say(ctx context.Context, text string) (string, error) {
	res, err := f.client.Parse(ctx, text)
	if err != nil {
		return "", xerrors.Errorf("fail to parse the text: %w", err)
	}

	if len(res.Results) == 0 || len(res.Results[0].Chunks) == 0 {
		return "", nil
	}

	specialtyID := -1

	chunkMap := make(map[int][]da.Chunk)
	for _, chunk := range res.Results[0].Chunks {
		chunkMap[chunk.Dependency] = append(chunkMap[chunk.Dependency], chunk)

		if specialtyID < 0 {
			for _, m := range chunk.Morphemes {
				if m.Surface == "得意" && m.POS == "名詞" {
					specialtyID = chunk.ID
					break
				}
			}
		}
	}

	if specialtyID < 0 {
		return "", nil
	}

	deps, ok := chunkMap[specialtyID]
	if !ok {
		return "", nil
	}

	var subject, specialty string
	for _, dep := range deps {
		mlen := len(dep.Morphemes)
		if mlen < 2 {
			continue
		}

		lastWord := dep.Morphemes[mlen-1]
		if lastWord.POS != "助詞" {
			continue
		}

		switch lastWord.Surface {
		case "は":
			subject = ""
			for i := 0; i < mlen-1; i++ {
				m := dep.Morphemes[i]
				subject += m.Surface
			}
		case "が", "も":
			specialty = ""
			for i := 0; i < mlen-1; i++ {
				m := dep.Morphemes[i]
				specialty += m.Surface
			}
		}
	}
	if specialty == "" {
		return "", nil
	}

	if sub, ok := firstPerson(subject); ok {
		subject = sub
	}

	return fmt.Sprintf("すごーい！%sは%sが得意なフレンズなんだね！", subject, specialty), nil
}

func firstPerson(subject string) (string, bool) {
	if subject == "" {
		return "きみ", true
	}

	plural := false
	if strings.HasSuffix(subject, "達") || strings.HasSuffix(subject, "たち") || strings.HasSuffix(subject, "ら") {
		chars := []rune(subject)
		subject = string(chars[0 : len(chars)-1])
		plural = true
	}

	for _, fp := range FirstPersons {
		if subject == fp {
			if plural {
				return "きみたち", true
			}
			return "きみ", true
		}
	}

	return "", false
}
