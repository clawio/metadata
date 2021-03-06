package simple

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/clawio/entities"
	"github.com/clawio/metadata/metadatacontroller"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var user = &entities.User{Username: "test"}

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
	require.Equal(suite.T(), "myblob", info.PathSpec)
	require.Equal(suite.T(), int64(1), info.Size)
	require.Equal(suite.T(), "", info.Checksum)
	require.Equal(suite.T(), "", info.MimeType)
	require.Equal(suite.T(), entities.ObjectTypeBLOB, info.Type)
}

func (suite *TestSuite) TestExamineObject_withNotFound() {
	_, err := suite.metadataController.ExamineObject(user, "notexists")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestListTree() {
	err := os.MkdirAll(suite.controller.getStoragePath(user, "testlisttree"), 0755)
	require.Nil(suite.T(), err)
	err = os.MkdirAll(suite.controller.getStoragePath(user, "testlisttree/othertree"), 0755)
	require.Nil(suite.T(), err)
	infos, err := suite.metadataController.ListTree(user, "testlisttree")
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), 1, len(infos))
}

func (suite *TestSuite) TestListTree_withNotFound() {
	_, err := suite.metadataController.ListTree(user, "notexists")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestDeleteObject() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = suite.metadataController.DeleteObject(user, "myblob")
	require.Nil(suite.T(), err)
}

func (suite *TestSuite) TestMoveBLOBObject() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "testmoveblobobject"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "testmoveblobobject", "othertestmoveblobobject")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestMoveTreeObject() {
	err := os.MkdirAll(suite.controller.getStoragePath(user, "testmovetree"), 0755)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "testmovetree", "othertestmovetree")
	require.Nil(suite.T(), err)
}

func (suite *TestSuite) TestMoveBLOBObject_overExistingBLOB() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob2"), []byte("2"), 0644)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "myblob", "myblob2")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestMoveBLOBObject_overExistingTree() {
	err := os.MkdirAll(suite.controller.getStoragePath(user, "mytree"), 0755)
	require.Nil(suite.T(), err)
	err = ioutil.WriteFile(suite.controller.getStoragePath(user, "mytree/myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "myblob", "mytree")
	require.NotNil(suite.T(), err)
	// err is the following
	// &os.LinkError{Op:"rename", Old:"/tmp/t/test/myblob", New:"/tmp/t/test/mytree", Err:0x15}
	// Err = "rename /tmp/t/test/myblob /tmp/t/test/mytree: is a directory"
}
func (suite *TestSuite) TestMoveTreeObject_overExistingBLOB() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "testmovetreeoverblobblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = os.MkdirAll(suite.controller.getStoragePath(user, "testmovetreeoverblobtree"), 0755)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "testmovetreeoverblobtree", "testmovetreeoverblobblob")
	require.NotNil(suite.T(), err)
	// err is the following
	// &os.LinkError{Op:"rename", Old:"/tmp/t/test/testmovetreeoverblobtree", New:"/tmp/t/test/testmovetreeoverblobblob", Err:0x14}
	// Err = "rename /tmp/t/test/testmovetreeoverblobtree /tmp/t/test/testmovetreeoverblobblob: not a directory"
}

func (suite *TestSuite) TestMoveTreeObject_overExistingTree() {
	err := os.MkdirAll(suite.controller.getStoragePath(user, "testmovetreeobjectmytreeovertree"), 0755)
	require.Nil(suite.T(), err)
	err = os.MkdirAll(suite.controller.getStoragePath(user, "testmovetreeobjectotheremptytree"), 0755)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "testmovetreeobjectmytreeovertre", "testmovetreeobjectotheremptytree")
	require.NotNil(suite.T(), err)
	// err is the following
	// &os.LinkError{Op:"rename", Old:"/tmp/t/test/mytreeovertree", New:"/tmp/t/test/otheremptytree", Err:0x42}
	// Err = rename /tmp/t/test/mytreeovertree /tmp/t/test/otheremptytree: directory not empty
}
func (suite *TestSuite) TestMoveObject_withTargetNotFound() {
	err := ioutil.WriteFile(suite.controller.getStoragePath(user, "myblob"), []byte("1"), 0644)
	require.Nil(suite.T(), err)
	err = suite.metadataController.MoveObject(user, "myblob", "notexists/otherblob")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestMoveObject_withSourceNotFound() {
	err := suite.metadataController.MoveObject(user, "notexists", "otherblob")
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
