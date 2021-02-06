package audio

import (
	"io"

	"github.com/go-logr/logr"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

// PlaybackPipelineOpts are options passed to the playback pipeline.
type playbackPipelineOpts struct {
	PulseServer, DeviceName, SourceFormat string
	SourceRate, SourceChannels            int
}

type pipelineReader struct {
	rPipe    *io.PipeReader
	wPipe    *io.PipeWriter
	pipeline *gst.Pipeline
}

func newPlaybackPipelineReader(log logr.Logger, errors chan error, opts *playbackPipelineOpts) (rdr io.ReadCloser, err error) {
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return
	}

	elements, err := gst.NewElementMany("pulsesrc", "cutter", "opusenc", "webmmux", "appsink")
	if err != nil {
		return
	}
	pulsesrc, cutter, opusenc, webmmux, appsink := elements[0], elements[1], elements[2], elements[3], elements[4]

	if err = pulsesrc.SetProperty("server", opts.PulseServer); err != nil {
		return
	}

	if err = pulsesrc.SetProperty("device", opts.DeviceName); err != nil {
		return
	}

	pulsecaps := newRawCaps(opts.SourceFormat, opts.SourceRate, opts.SourceChannels)

	r, w := io.Pipe()

	app.SinkFromElement(appsink).SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(self *app.Sink) gst.FlowReturn {
			sample := self.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}
			if _, err := io.Copy(w, buffer.Reader()); err != nil {
				return gst.FlowError
			}
			return gst.FlowOK
		},
	})

	if err = pipeline.AddMany(elements...); err != nil {
		return
	}
	if err = pulsesrc.LinkFiltered(cutter, pulsecaps); err != nil {
		return
	}
	if err = gst.ElementLinkMany(cutter, opusenc, webmmux, appsink); err != nil {
		return
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageError:
			log.Error(err, "Error from pipeline")
			errors <- msg.ParseError()
		case gst.MessageEOS:
			log.Info("Pipeline has reached EOS")
			errors <- app.ErrEOS
		case gst.MessageElement:
		default:
			log.Info(msg.String())
		}
		return true
	})

	if err = pipeline.SetState(gst.StatePlaying); err != nil {
		return
	}

	return &pipelineReader{
		rPipe:    r,
		wPipe:    w,
		pipeline: pipeline,
	}, nil
}

func (r *pipelineReader) Read(p []byte) (int, error) {
	return r.rPipe.Read(p)
}

func (r *pipelineReader) Close() error {
	if err := r.pipeline.BlockSetState(gst.StateNull); err != nil {
		return err
	}
	return r.wPipe.Close()
}
