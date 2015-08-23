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

import twodee "../lib/twodee"

type AudioSystem struct {
	app                   *Application
	bgm                   *twodee.Music
	bgmObserverId         int
	pauseMusicObserverId  int
	resumeMusicObserverId int
	musicToggle           int32
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

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.bgmObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.bgm.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		bgm *twodee.Music
	)

	if bgm, err = twodee.NewMusic("resources/music/bgm1.ogg"); err != nil {
		return
	}

	audioSystem = &AudioSystem{
		app:         app,
		bgm:         bgm,
		musicToggle: 1,
	}
	audioSystem.bgmObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	return
}
