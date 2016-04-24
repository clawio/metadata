package simple

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/clawio/entities"
	"github.com/clawio/entities/mocks"
	"github.com/clawio/metadata/metadatacontroller"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var user = &mocks.MockUser{Username: "test"}

type TestSuite struct {
	suite.Suite
	metadataController metadatacontroller.MetaDataController
	controller         *controller
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (suite *TestSuite) SetupTest() {
	opts := &Options{
		MetaDataDir: "/tmp",
		TempDir:     "/tmp",
	}
	metadataController := New(opts)
	// create homedir for user test
	err := os.MkdirAll("/tmp/t/test", 0755)
	require.Nil(suite.T(), err)
	suite.metadataController = metadataController
	suite.controller = suite.metadataController.(*controller)

	// configure user mock
	user.On("GetUsername").Return("test")
}
func (suite *TestSuite) TeardownTest() {
	os.RemoveAll("/tmp/t")
}
func (suite *TestSuite) New() {
	opts := &Options{
		MetaDataDir: "/tmp",
		TempDir:     "/tmp",
	}
	require.IsType(suite.T(), &controller{}, New(opts))
}
func (suite *TestSuite) TestNew_withNilOptions() {
	require.IsType(suite.T(), &controller{}, New(nil))
}
func (suite *TestSuite) TestInit() {
	err := suite.metadataController.Init(user)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestInit_withError() {
	suite.controller.metaDataDir = "/i/cannot/write/here"
	err := suite.metadataController.Init(user)
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestExamineObject() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	info, err := suite.metadataController.ExamineObject(user, "myblob")
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), "myblob", info.GetPathSpec())
	require.Equal(suite.T(), uint64(1), info.GetSize())
	require.Equal(suite.T(), "", info.GetChecksum())
	require.Equal(suite.T(), "", info.GetMimeType())
	require.Equal(suite.T(), entities.ObjectTypeBLOB, info.GetType())
}

func (suite *TestSuite) TestExamineObject_withNotFound() {
	_, err := suite.metadataController.ExamineObject(user, "notexists")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestListTree() {
	err := os.MkdirAll(suite.controller.getStoragePath(user, "mytree"), 0755)
	require.Nil(suite.T(), err)
	err = os.MkdirAll(suite.controller.getStoragePath(user, "mytree/othertree"), 0755)
	require.Nil(suite.T(), err)
	infos, err := suite.metadataController.ListTree(user, "mytree")
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), 1, len(infos))
}

func (suite *TestSuite) TestListTree_withNotFound() {
	_, err := suite.metadataController.ListTree(user, "notexists")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestListTree_withBLOB() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	_, err = suite.metadataController.ListTree(user, "myblob")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestgetMimeType() {
	mime := suite.controller.getMimeType("", entities.ObjectTypeTree)
	require.Equal(suite.T(), entities.ObjectTypeTreeMimeType, mime)
}

func (suite *TestSuite) TestgetMimeType_pdf() {
	mime := suite.controller.getMimeType("myblob.pdf", entities.ObjectTypeBLOB)
	require.Equal(suite.T(), "application/pdf", mime)
}

/*
func (suite *TestSuite) TestUpload() {
	reader := strings.NewReader("1")
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestUpload_withBadTempDir() {
	suite.controller.tempDir = "/this/does/not/exist"
	reader := strings.NewReader("1")
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestUpload_withChecksum() {
	suite.controller.checksum = "md5"
	reader := strings.NewReader("1")
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestUpload_withWrongChecksum() {
	suite.controller.checksum = "xyz"
	reader := strings.NewReader("1")
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestUpload_withClientChecksum() {
	suite.controller.checksum = "md5"
	suite.controller.verifyClientChecksum = true
	reader := strings.NewReader("1")
	// md5 checksum of 1 is c4ca4238a0b923820dcc509a6f75849b
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "md5:c4ca4238a0b923820dcc509a6f75849b")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestUpload_withWrongClientChecksum() {
	suite.controller.checksum = "md5"
	suite.controller.verifyClientChecksum = true
	reader := strings.NewReader("1")
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "md5:")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestUpload_withBadMetaDataDir() {
	suite.controller.metadataDir = "/this/does/not/exist"
	reader := strings.NewReader("1")
	err := suite.metadataController.UploadBLOB(user, "myblob", reader, "")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestDownload() {
	p := path.Join(suite.controller.tempDir, "t", "test", "myblob")
	err := ioutil.WriteFile(p, []byte("1"), 0644)
	require.Nil(suite.T(), err)
	reader, err := suite.metadataController.DownloadBLOB(user, "myblob")
	require.Nil(suite.T(), err)
	metadata, err := ioutil.ReadAll(reader)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), "1", string(metadata))
}
func (suite *TestSuite) TestDownload_withBadMetaDataDir() {
	suite.controller.metadataDir = "/this/does/not/exist"
	_, err := suite.metadataController.DownloadBLOB(user, "myblob")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcomputeChecksum_withBadChecksum() {
	suite.controller.checksum = "xyz"
	p := path.Join(suite.controller.tempDir, "t", "test", "myblob")
	err := ioutil.WriteFile(p, []byte("1"), 0644)
	require.Nil(suite.T(), err)
	_, err = suite.controller.computeChecksum(p)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcomputeChecksum_withNoFile() {
	suite.controller.checksum = "md5"
	_, err := suite.controller.computeChecksum("/this/does/not/exist/myblob")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcomputeChecksum_md5() {
	suite.controller.checksum = "md5"
	p := path.Join(suite.controller.tempDir, "t", "test", "myblob")
	err := ioutil.WriteFile(p, []byte("1"), 0644)
	require.Nil(suite.T(), err)
	checksum, err := suite.controller.computeChecksum(p)
	require.Nil(suite.T(), err)

	// md5 checksum of "1" is c4ca4238a0b923820dcc509a6f75849b
	require.Equal(suite.T(), "md5:c4ca4238a0b923820dcc509a6f75849b", checksum)
}
func (suite *TestSuite) TestcomputeChecksum_adler32() {
	suite.controller.checksum = "adler32"
	p := path.Join(suite.controller.tempDir, "t", "test", "myblob")
	err := ioutil.WriteFile(p, []byte("1"), 0644)
	require.Nil(suite.T(), err)
	checksum, err := suite.controller.computeChecksum(p)
	require.Nil(suite.T(), err)

	// adler32 checksum of "1" is 00320032
	require.Equal(suite.T(), "adler32:00320032", checksum)
}
func (suite *TestSuite) TestcomputeChecksum_sha1() {
	suite.controller.checksum = "sha1"
	p := path.Join(suite.controller.tempDir, "t", "test", "myblob")
	err := ioutil.WriteFile(p, []byte("1"), 0644)
	require.Nil(suite.T(), err)
	checksum, err := suite.controller.computeChecksum(p)
	require.Nil(suite.T(), err)

	// sha1 checksum of "1" is 356a192b7913b04c54574d18c28d46e6395428ab
	require.Equal(suite.T(), "sha1:356a192b7913b04c54574d18c28d46e6395428ab", checksum)
}
func (suite *TestSuite) TestcomputeChecksum_sha256() {
	suite.controller.checksum = "sha256"
	p := path.Join(suite.controller.tempDir, "t", "test", "myblob")
	err := ioutil.WriteFile(p, []byte("1"), 0644)
	require.Nil(suite.T(), err)
	checksum, err := suite.controller.computeChecksum(p)
	require.Nil(suite.T(), err)

	// sha256 checksum of "1" is 6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b
	require.Equal(suite.T(), "sha256:6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", checksum)
}

*/
func (suite *TestSuite) TestsecureJoin() {
	paths := []struct {
		given    []string
		expected string
	}{
		{
			[]string{"relativePath/t/test"},
			"relativePath/t/test",
		},
		{
			[]string{"../../relativePath/t/test"},
			"../../relativePath/t/test",
		},
		{
			[]string{"../../relativePath/t/test", "../../../../"},
			"../../relativePath/t/test",
		},
		{
			[]string{"/abspath/t/test"},
			"/abspath/t/test",
		},
		{
			[]string{"/abspath/t/test", "../../.."},
			"/abspath/t/test",
		},
	}

	for _, v := range paths {
		require.Equal(suite.T(), v.expected, secureJoin(v.given...))
	}
}
