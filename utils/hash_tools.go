package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
	"hash/crc32"
	"io"
	"os"

	"github.com/stratosnet/sds/utils/crypto"
)

// CalcCRC32
func CalcCRC32(data []byte) uint32 {
	iEEE := crc32.NewIEEE()
	io.WriteString(iEEE, string(data))
	return iEEE.Sum32()
}

// CalcFileMD5
func CalcFileMD5(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		Log(err.Error())
		return nil
	}
	defer file.Close()
	MD5 := md5.New()
	io.Copy(MD5, file)
	return MD5.Sum(nil)
}

// CalcFileCRC32
func CalcFileCRC32(filePath string) uint32 {
	file, err := os.Open(filePath)
	if err != nil {
		Log(err.Error())
		return 0
	}
	defer file.Close()
	iEEE := crc32.NewIEEE()
	io.Copy(iEEE, file)
	return iEEE.Sum32()
}

// CalcFileHash
// @notice keccak256(md5(file))
func CalcFileHash(filePath, encryptionTag string) string {
	if filePath == "" {
		Log(errors.New("CalcFileHash: missing file path"))
		return ""
	}
	data := append([]byte(encryptionTag), CalcFileMD5(filePath)...)
	return calcFileHash(data)
}

// CalcHash
func CalcHash(data []byte) string {
	return hex.EncodeToString(crypto.Keccak256(data))
}

// CalcHash
func CalcSliceHash(data []byte, fileHash string) string {
	fileCid, _ := cid.Decode(fileHash)
	fileKeccak256 := fileCid.Hash()
	sliceKeccak256, _ := mh.Sum(data, mh.KECCAK_256, 20)
	if len(fileKeccak256) != len(sliceKeccak256) {
		Log(errors.New("length of fileKeccak256 and sliceKeccak256 doesn't match"))
		return ""
	}
	sliceHash := make([]byte, len(fileKeccak256))
	for i := 0; i < len(fileKeccak256); i++ {
		sliceHash[i] = fileKeccak256[i] ^ sliceKeccak256[i]
	}
	sliceHash, _ = mh.Sum(sliceHash, mh.KECCAK_256, 20)
	sliceCid := cid.NewCidV1(cid.Raw, sliceHash)
	encoder, _ := mbase.NewEncoder(mbase.Base32hex)
	return sliceCid.Encode(encoder)
}

func calcFileHash(data []byte) string {
	fileHash, _ := mh.Sum(data, mh.KECCAK_256, 20)
	fileCid := cid.NewCidV1(cid.Raw, fileHash)
	encoder, _ := mbase.NewEncoder(mbase.Base32hex)
	return fileCid.Encode(encoder)
}
