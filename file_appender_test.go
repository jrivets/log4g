package log4g

import (
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type faConfigSuite struct {
}

var _ = Suite(&faConfigSuite{})

func (s *faConfigSuite) TestFileAppenderName(c *C) {
	c.Assert(faFactory.Name(), Equals, "log4g/fileAppender")
}

func (s *faConfigSuite) TestNewAppender(c *C) {
	app, err := faFactory.NewAppender(nil)
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %ee"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %ee", "fileName": "fn"}) //bad layout
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "-1"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "abc"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true", "maxFileSize": "10"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxDiskSpace": "10"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxDiskSpace": "2G", "rotate": "daily2"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxDiskSpace": "2K", "rotate": "daily"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxDiskSpace": "2M", "rotate": "size"})
	c.Assert(app, NotNil)
	c.Assert(err, IsNil)
	c.Assert(app.(*fileAppender).rotate, Equals, rsSize)
	app.Shutdown()
}

func (s *faConfigSuite) TestShutdown(c *C) {
	app, err := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000",
		"maxFileSize": "2K", "maxDiskSpace": "2M", "rotate": "daily"})
	c.Assert(app, NotNil)
	c.Assert(err, IsNil)
	fa := app.(*fileAppender)
	c.Assert(fa.file, IsNil)
	c.Assert(fa.fileAppend, Equals, true)
	c.Assert(fa.fileName, Equals, "fn")
	c.Assert(fa.maxSize, Equals, int64(2000))
	c.Assert(fa.maxDiskSpace, Equals, int64(2000000))
	c.Assert(fa.rotate, Equals, rsDaily)
	ok := false
	select {
	case _, ok = <-fa.msgChannel:
		ok = true
		break
	default:
	}
	c.Assert(ok, Equals, false)
	app.Shutdown()

	_, ok = <-fa.msgChannel
	c.Assert(ok, Equals, false)
}

func (s *faConfigSuite) TestAppendDiskSpace(c *C) {
	defer removeFiles("456____test____log___file")
	fa := writeLogs(c, map[string]string{"layout": "%p", "fileName": "456____test____log___file", "buffer": "1000",
		"maxFileSize": "1K", "maxDiskSpace": "10K", "rotate": "size"}, 5000)
	c.Check(fa.stat.chunksSize >= 9000 && fa.stat.chunksSize <= 10000, Equals, true)
}

func (s *faConfigSuite) TestAppendToExistingOne(c *C) {
	defer removeFiles("____test____log___file")
	params := map[string]string{"layout": "%p", "fileName": "____test____log___file", "buffer": "1000",
		"maxFileSize": "2Gib", "rotate": "none", "append": "true"}
	fa := writeLogs(c, params, 10000)
	size := fa.stat.size
	fa = writeLogs(c, params, 10000)
	c.Check(fa.stat.size, Equals, size*2)
	params["append"] = "false"
	fa = writeLogs(c, params, 10000)
	c.Check(fa.stat.size, Equals, size)
}

func (s *faConfigSuite) TestSizeRotation(c *C) {
	app, _ := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000",
		"maxFileSize": "2K", "masDiskSpace": "10K", "rotate": "daily"})
	fa := app.(*fileAppender)
	fa.stat.size = 2000
	c.Assert(fa.sizeRotation(), Equals, false)
	fa.stat.size++
	c.Assert(fa.sizeRotation(), Equals, true)
	app.Shutdown()
}

func (s *faConfigSuite) TestTimeRotation(c *C) {
	app, _ := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000",
		"maxFileSize": "2K", "maxLines": "2Mib", "rotate": "daily"})
	fa := app.(*fileAppender)
	fa.stat.startTime = time.Now()
	c.Assert(fa.timeRotation(), Equals, false)

	fa.stat.startTime = time.Now().AddDate(0, 0, -1)
	c.Assert(fa.timeRotation(), Equals, true)

	fa.stat.startTime = time.Now().AddDate(0, 0, -7)
	c.Assert(fa.timeRotation(), Equals, true)

	fa.stat.startTime = time.Now().AddDate(0, 0, 1)
	c.Assert(fa.timeRotation(), Equals, true)
	app.Shutdown()
}

func (s *faConfigSuite) TestisRotationNeededNone(c *C) {
	defer removeFiles("____test____log___file2")
	app, _ := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "____test____log___file2", "buffer": "1000",
		"maxFileSize": "2K", "maxLines": "2K"})
	fa := app.(*fileAppender)
	c.Check(fa.isRotationNeeded(), Equals, true)

	app.Append(&Event{INFO, time.Now(), "abc", "def"})
	for fa.file == nil {
		time.Sleep(time.Millisecond)
	}

	c.Check(fa.isRotationNeeded(), Equals, false)
	fa.stat.size = 2001
	c.Check(fa.isRotationNeeded(), Equals, false)
	fa.stat.startTime = time.Now().AddDate(0, 0, -3)
	c.Check(fa.isRotationNeeded(), Equals, false)

	app.Shutdown()
}

func (s *faConfigSuite) TestisRotationNeededSize(c *C) {
	defer removeFiles("____test____log___file2")
	app, _ := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "____test____log___file2", "buffer": "1000",
		"maxFileSize": "2K", "maxLines": "2K", "rotate": "size"})
	fa := app.(*fileAppender)
	c.Check(fa.isRotationNeeded(), Equals, true)

	app.Append(&Event{INFO, time.Now(), "abc", "def"})
	for fa.file == nil {
		time.Sleep(time.Millisecond)
	}

	c.Check(fa.isRotationNeeded(), Equals, false)
	fa.stat.startTime = time.Now().AddDate(0, 0, -3)
	c.Check(fa.isRotationNeeded(), Equals, false)
	fa.stat.size = 2001
	c.Check(fa.isRotationNeeded(), Equals, true)

	app.Shutdown()
}

func (s *faConfigSuite) TestisRotationNeededDaily(c *C) {
	defer removeFiles("____test____log___file2")
	app, _ := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "____test____log___file2", "buffer": "1000",
		"maxFileSize": "2K", "maxLines": "2K", "rotate": "daily"})
	fa := app.(*fileAppender)
	c.Check(fa.isRotationNeeded(), Equals, true)

	app.Append(&Event{INFO, time.Now(), "abc", "def"})
	for fa.file == nil {
		time.Sleep(time.Millisecond)
	}

	c.Check(fa.isRotationNeeded(), Equals, false)
	fa.stat.startTime = time.Now().AddDate(0, 0, -3)
	c.Check(fa.isRotationNeeded(), Equals, true)
	fa.stat.startTime = time.Now()
	fa.stat.size = 2001
	c.Check(fa.isRotationNeeded(), Equals, true)
	app.Shutdown()
}

func removeFiles(prefix string) {
	archiveName, _ := filepath.Abs(prefix)
	dir := filepath.Dir(archiveName)
	baseName := filepath.Base(archiveName)
	fileInfos, _ := ioutil.ReadDir(dir)
	for _, fInfo := range fileInfos {
		if fInfo.IsDir() || !strings.HasPrefix(fInfo.Name(), baseName) {
			continue
		}
		os.Remove(fInfo.Name())
	}
}

func writeLogs(c *C, params map[string]string, count int) *fileAppender {
	app, _ := faFactory.NewAppender(params)
	c.Assert(app, NotNil)
	fa := app.(*fileAppender)
	for idx := 0; idx < count; idx++ {
		app.Append(&Event{INFO, time.Now(), "abc", "def"})
	}
	app.Shutdown()
	return fa
}
