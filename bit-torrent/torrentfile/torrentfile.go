package torrentfile

import (
	"bytes"
	"fmt"
	"io"
	//"crypto/rand"
	"crypto/sha1"

	"github.com/jackpal/bencode-go"
)

// TorrentFile encodes the metadata from a .torrent file
// This means it hold ionformation about the .torrent file
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

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

// Implement the helper functions here
// @TODO : 
// Hash , splitPieceHashes

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile , error) {
	infohash , err := bto.Info.hash()
	if err != nil {
		return TorrentFile{} , err
	}
	pieceHashes , err := bto.Info.splitPieceHash()
	if err != nil {
		return TorrentFile{} , err
	}
	t := TorrentFile {
		Announce: bto.Announce,
		InfoHash: infohash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PiecesLength,
		Length: bto.Info.Length,
		Name: bto.Info.Name,
	}
	return t , nil
}

// splitPieceHash function
func (i *bencodeInfo) splitPieceHash() ([][20]byte , error){
	hashLen := 20 // Length of hash
	buf := []byte(i.Pieces)
	if len(buf) % hashLen != 0 {
		err := fmt.Errorf("Received wrong pieces of length %d" , len(buf))
		return nil , err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte , numHashes)

	for i := 0 ; i < numHashes ; i++{
		copy(hashes[i][:] , buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes , nil
}

// Hash function
func (i *bencodeInfo) hash() ([20]byte , error){
	var buf bytes.Buffer
	err := bencode.Marshal(&buf , *i)
	if err != nil {
		return [20]byte{} , err
	}
	h := sha1.Sum(buf.Bytes())
	return h , nil
}