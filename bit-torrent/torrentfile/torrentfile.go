package torrentfile

import("io"
	"bytes"
	"github.com/jackpal/bencode-go"


)

type bencodeInfo struct {
	Pieces string `bencode:"peices"`
	PiecesLength int `bencode:"peices length"`
	Length int `bencode:"length"`
	Name string `bencode:"name"`
}

type bencodeTorrent struct {
	Info bencodeInfo `bencode:"info"`
	Announce string `bencode:"announce"`
}

// Open parses a torrent file
func Open(r io.Reader) (*bencodeTorrent,error){
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r , &bto)
	if err != nil {
		return nil , err
	}
	return &bto , nil
}
