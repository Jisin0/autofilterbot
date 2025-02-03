// Package fileid is cherry picked funtions and types extracted from github.com/gotd/td to encode and decode file ids.
package fileid

import (
	"encoding/base64"
)

// FileID represents parsed Telegram Bot API file_id.
type FileID struct {
	Type          Type
	DC            int
	ID            int64
	AccessHash    int64
	FileReference []byte
	URL           string
}

const (
	webLocationFlag        = 1 << 24
	fileReferenceFlag      = 1 << 25
	latestSubVersion       = 34
	persistentIDVersionOld = 2
	persistentIDVersionMap = 3
	persistentIDVersion    = 4
)

func base64Encode(s []byte) string {
	return base64.RawURLEncoding.EncodeToString(s)
}

func rleEncode(s []byte) (r []byte) {
	var count byte
	for _, cur := range s {
		if cur == 0 {
			count++
			continue
		}

		if count > 0 {
			r = append(r, 0, count)
			count = 0
		}
		r = append(r, cur)
	}
	if count > 0 {
		r = append(r, 0, count)
	}

	return r
}

func (f *FileID) encodeLatestFileID(b *Buffer) {
	hasWebLocation := f.URL != ""
	hasReference := len(f.FileReference) != 0

	{
		typeID := f.Type
		if hasWebLocation {
			typeID |= webLocationFlag
		}
		if hasReference {
			typeID |= fileReferenceFlag
		}
		b.PutUint32(uint32(typeID))
	}
	b.PutUint32(uint32(f.DC))
	if hasReference {
		b.PutBytes(f.FileReference)
	}
	if hasWebLocation {
		b.PutString(f.URL)
		return
	}
	b.PutLong(f.ID)
	b.PutLong(f.AccessHash)

	b.Buf = append(b.Buf, latestSubVersion)
}

// EncodeFileID parses FileID to a string.
func EncodeFileID(id FileID) (string, error) {
	var buf Buffer
	id.encodeLatestFileID(&buf)
	buf.Buf = append(buf.Buf, persistentIDVersion)
	buf.Buf = rleEncode(buf.Buf)
	return base64Encode(buf.Buf), nil
}
