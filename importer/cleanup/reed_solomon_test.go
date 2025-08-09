package cleanup

//
//import (
//	"github.com/stretchr/testify/suite"
//	"testing"
//)
//
//type ReedSolomonSuite struct {
//	suite.Suite
//}
//
//// In order for 'go test' to run this suite, we need to create
//// a normal test function and pass our suite to suite.Run
//func TestReedSolomon(t *testing.T) {
//	suite.Run(t, new(ReedSolomonSuite))
//}
//
//func (suite *ReedSolomonSuite) TestBestChunk() {
//	//stats := BestChunking(1, RS9x3)
//	//suite.Require().Zero(stats)
//	stats := BestChunking(int64(c64k), RS9x3)
//	suite.Require().Zero(stats)
//	stats = BestChunking(int64(c64k)+1, RS9x3)
//	suite.Require().Zero(stats)
//	stats = BestChunking(int64(c64k)*9, RS9x3)
//	suite.Require().Zero(stats)
//	stats = BestChunking(int64(c64k)*9 +1, RS9x3)
//	suite.Require().Equal(stats.ChunkSize, c64k,"Chunk Size")
//	suite.Require().Equal(stats.Data, 10,"Data")
//	suite.Require().Equal(stats.Empty, 8,"Empty")
//	suite.Require().Equal(stats.Parity, 6,"Parity")
//
//}
