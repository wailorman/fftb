package ff

import "github.com/wailorman/fftb/pkg/files"

// Progress _
type Progress struct {
	framesProcessed string
	currentTime     string
	currentBitrate  string
	progress        float64
	speed           string
	fps             float64
	file            files.Filer
}

// Progressable _
type Progressable interface {
	FramesProcessed() string
	CurrentTime() string
	CurrentBitrate() string
	Progress() float64
	Speed() string
	FPS() float64
	File() files.Filer
}

// FramesProcessed _
func (p *Progress) FramesProcessed() string {
	return p.framesProcessed
}

// CurrentTime _
func (p *Progress) CurrentTime() string {
	return p.currentTime
}

// CurrentBitrate _
func (p *Progress) CurrentBitrate() string {
	return p.currentBitrate
}

// Progress _
func (p *Progress) Progress() float64 {
	return p.progress
}

// Speed _
func (p *Progress) Speed() string {
	return p.speed
}

// FPS _
func (p *Progress) FPS() float64 {
	return p.fps
}

// File _
func (p *Progress) File() files.Filer {
	return p.file
}
