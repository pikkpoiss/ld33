// Copyright 2015 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"../lib/twodee"
	"time"
)

const (
	MR_BONES_EFFECT_DURATION float64 = 2
	SPIKES_EFFECT_DURATION   float64 = 2
)

type AudioSystem struct {
	app                        *Application
	bgm                        *twodee.Music
	placeBlockEffect           *twodee.SoundEffect
	mrbonesEffect              *twodee.SoundEffect
	spikesEffect               *twodee.SoundEffect
	deathEffect                *twodee.SoundEffect
	bgmObserverId              int
	pauseMusicObserverId       int
	resumeMusicObserverId      int
	placeBlockEffectObserverId int
	mrbonesEffectObserverId    int
	spikesEffectObserverId     int
	deathEffectObserverId      int
	musicToggle                int32
	mrbonesLastPlayed          time.Time
	spikesLastPlayed           time.Time
}

func (a *AudioSystem) PlayBackgroundMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.bgm.Play(-1)
	}
}

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPaused() {
			twodee.ResumeMusic()
		}
	}
}

func (a *AudioSystem) PlayPlaceBlockEffect(e twodee.GETyper) {
	if a.placeBlockEffect.IsPlaying(2) == 0 {
		a.placeBlockEffect.PlayChannel(2, 1)
	}
}

func (a *AudioSystem) PlayMrBonesEffect(e twodee.GETyper) {
	if time.Since(a.mrbonesLastPlayed).Seconds() > MR_BONES_EFFECT_DURATION {
		a.mrbonesEffect.PlayChannel(3, 1)
		a.mrbonesLastPlayed = time.Now()
	}
}

func (a *AudioSystem) PlaySpikesEffect(e twodee.GETyper) {
	if time.Since(a.spikesLastPlayed).Seconds() > SPIKES_EFFECT_DURATION {
		a.spikesEffect.PlayChannel(4, 1)
		a.spikesLastPlayed = time.Now()
	}
}

func (a *AudioSystem) PlayDeathEffect(e twodee.GETyper) {
	if a.deathEffect.IsPlaying(5) == 0 {
		a.deathEffect.PlayChannel(5, 1)
	}
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.bgmObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayPlaceBlockEffect, a.placeBlockEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayMrBonesEffect, a.mrbonesEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlaySpikesEffect, a.spikesEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayDeathEffect, a.deathEffectObserverId)
	a.bgm.Delete()
	a.placeBlockEffect.Delete()
	a.mrbonesEffect.Delete()
	a.spikesEffect.Delete()
	a.deathEffect.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		bgm              *twodee.Music
		placeBlockEffect *twodee.SoundEffect
		mrbonesEffect    *twodee.SoundEffect
		spikesEffect     *twodee.SoundEffect
		deathEffect      *twodee.SoundEffect
	)

	if bgm, err = twodee.NewMusic("resources/music/bgm1.ogg"); err != nil {
		return
	}

	if placeBlockEffect, err = twodee.NewSoundEffect("resources/music/place-block.ogg"); err != nil {
		return
	}

	if mrbonesEffect, err = twodee.NewSoundEffect("resources/music/deep-laugh.ogg"); err != nil {
		return
	}

	if spikesEffect, err = twodee.NewSoundEffect("resources/music/spikes.ogg"); err != nil {
		return
	}

	if deathEffect, err = twodee.NewSoundEffect("resources/music/no.ogg"); err != nil {
		return
	}

	mrbonesLastPlayed := time.Now()
	spikesLastPlayed := time.Now()

	audioSystem = &AudioSystem{
		app:               app,
		bgm:               bgm,
		placeBlockEffect:  placeBlockEffect,
		mrbonesEffect:     mrbonesEffect,
		spikesEffect:      spikesEffect,
		deathEffect:       deathEffect,
		musicToggle:       1,
		mrbonesLastPlayed: mrbonesLastPlayed,
		spikesLastPlayed:  spikesLastPlayed,
	}
	audioSystem.bgmObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	audioSystem.placeBlockEffectObserverId = app.GameEventHandler.AddObserver(PlayPlaceBlockEffect, audioSystem.PlayPlaceBlockEffect)
	audioSystem.mrbonesEffectObserverId = app.GameEventHandler.AddObserver(PlayMrBonesEffect, audioSystem.PlayMrBonesEffect)
	audioSystem.spikesEffectObserverId = app.GameEventHandler.AddObserver(PlaySpikesEffect, audioSystem.PlaySpikesEffect)
	audioSystem.deathEffectObserverId = app.GameEventHandler.AddObserver(PlayDeathEffect, audioSystem.PlayDeathEffect)
	return
}
