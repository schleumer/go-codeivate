// Copyright AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA

// +build ignore

package main

import (
  ui "github.com/gizak/termui"
  "time"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strings"
  "sort"
  "fmt"
  "math"
  "flag"
  "log"
)

type LangStatistic struct {
  Level string `json:"level"`
  Points float64 `json:"points"`
}

type Language struct {
  Name string
  Points float64
  Level int
  Percent float64
}

type PlatformStatistic struct {
  PercentWork float64 `json:"percent_work"`
  Points float64 `json:"points"`
  Time float64 `json:"time"`
}

type Platform struct {
  Name string
  PercentWork float64
  Points float64
  Time float64
}

type UserStatistic struct {
  CurrentLanguage string `json:"current_language"`
  FocusLevel string `json:"focus_level"`
  FocusPoints float64 `json:"focus_points"`
  Level string `json:"level"`
  MaxStreak float64 `json:"max_streak"`
  Name string `json:"name"`
  ProgrammingNow bool `json:"programming_now"`
  StreakingNow bool `json:"streaking_now"`
  TimeSpent float64 `json:"time_spent"`
  TotalDaysCoded float64 `json:"total_days_coded"`
  TotalFlowStates float64 `json:"total_flow_states"`
  Platforms map[string]PlatformStatistic `json:"platforms"`
  Languages map[string]LangStatistic `json:"languages"`
}

type ByLevel []Language
func (v ByLevel) Len() int { return len(v) }
func (v ByLevel) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v ByLevel) Less(i, j int) bool { 
  if v[i].Level == v[j].Level {
    if v[i].Points == v[j].Points {
      return v[i].Name > v[j].Name
    } else {
      return v[i].Points > v[j].Points    
    }
    
  } else {
    return v[i].Level > v[j].Level 
  }
}

type PlatformByPoints []Platform
func (v PlatformByPoints) Len() int { return len(v) }
func (v PlatformByPoints) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v PlatformByPoints) Less(i, j int) bool { 
  if v[i].Points == v[j].Points {
    if v[i].Time == v[j].Time {
      return v[i].Name > v[j].Name
    } else {
      return v[i].Time > v[j].Time    
    }
  } else {
    return v[i].Points > v[j].Points 
  }
}

type Level struct {
  Number int
  Percent float64
}

func ParseLevel(strLevel string) (Level, error) {
  splited := strings.Split(strLevel, ".")
  level, err := strconv.Atoi(splited[0])
  if err != nil {
    return Level{}, err
  }
  percent, err := strconv.ParseFloat(splited[1], 64)
  if err != nil {
    return Level{}, err
  }
  return Level{level, percent}, nil
}

func boot() string {
  // yeah, if there's no username you'll see my profile :3
  var username string
  flag.StringVar(&username, "username", "schleumer", "your codeivate username")
  flag.Parse()

  err := ui.Init()
  if err != nil {
    panic(err)
  }
  defer ui.Close()

  ui.UseTheme("helloworld")
  
  done := make(chan bool)
  redraw := make(chan bool)
  error := make(chan string)

  ui.Body.Align()

  update := func () {
    for {
      // restart body
      ui.Body = ui.NewGrid()
      ui.Body.Width = ui.TermWidth()

      resp, err := http.Get(fmt.Sprintf("http://codeivate.com/users/%s.json", username))
      if err != nil {
        error <- "Error on request"
        return
      }

      defer resp.Body.Close()

      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        error <- "Error on reading response"
        return
      }

      var statistic UserStatistic
      err = json.Unmarshal(body, &statistic)
      if err != nil {
        error <- fmt.Sprintf("Error on unmarshaling, probably the user %s doesn't exists", username)
        return
      }

      userLevel, err := ParseLevel(statistic.Level)
      if err != nil {
        println(err)
      }

      // i have no idea what i'm doing, but it's fucking hardcore
      hours := math.Floor(statistic.TimeSpent / 3600)
      minutes := math.Floor((statistic.TimeSpent - (hours * 3600)) / 60)

      parUser := ui.NewPar(fmt.Sprintf("Level: %d - Percent: %.0f - Time: %.0f hours %.0f minutes - Current Language: %s", userLevel.Number, userLevel.Percent, hours, minutes, statistic.CurrentLanguage))
      parUser.Height = 3
      parUser.Width = 50
      parUser.TextFgColor = ui.ColorWhite
      parUser.Border.Label = "User Info"
      parUser.Border.FgColor = ui.ColorCyan

      ui.Body.AddRows(
          ui.NewRow(ui.NewCol(12, 0, parUser)))


      var platformsContent []string

      var platforms []Platform

      for name, platform := range statistic.Platforms {
        platforms = append(platforms, Platform{name, platform.PercentWork, platform.Points, platform.Time})
      }

      sort.Sort(PlatformByPoints(platforms))

      for _, platform := range platforms {
        hours := math.Floor(platform.Time / 3600)
        minutes := math.Floor((platform.Time - (hours * 3600)) / 60)
        platformsContent = append(platformsContent, fmt.Sprintf("%s - Percent: %.2f - Points: %.2f - Time: %.0f hours %.0f minutes", platform.Name, platform.PercentWork, platform.Points, hours, minutes))
      }

      parWorkspace := ui.NewPar(strings.Join(platformsContent, "\n"))
      parWorkspace.Height = len(statistic.Platforms) + 2
      parWorkspace.Width = 50
      parWorkspace.TextFgColor = ui.ColorWhite
      parWorkspace.Border.Label = "Workspace Info"
      parWorkspace.Border.FgColor = ui.ColorCyan

      ui.Body.AddRows(
          ui.NewRow(ui.NewCol(12, 0, parWorkspace)))

      var languages []Language

      for name, lang := range statistic.Languages {
        level, err := ParseLevel(lang.Level)
        if err != nil {
          // don't give a fuck 
          continue
        }
        languages = append(languages, Language{name, lang.Points, level.Number, level.Percent})
      }

      if err != nil {
        // don't give a fuck 
        continue
      }

      sort.Sort(ByLevel(languages))
      
      for _, lang := range languages[:8] {
        g := ui.NewGauge()
        g.Percent = int(lang.Percent)
        g.Width = 50
        g.Height = 3
        g.Border.Label = fmt.Sprintf("%s Level: %d - Points: %.0f", lang.Name, lang.Level, lang.Points)
        g.BarColor = ui.ColorRed
        g.Border.FgColor = ui.ColorWhite
        g.Border.LabelFgColor = ui.ColorCyan
  
        ui.Body.AddRows(
          ui.NewRow(ui.NewCol(12, 0, g)))
      }

      ui.Render(ui.Body)
      ui.Body.Align()

      redraw <- true
      time.Sleep(time.Second * 10)
    }
  }
  

  evt := ui.EventCh()

  ui.Render(ui.Body)
  go update()

  for {
    select {
      case e := <-evt:
        if e.Type == ui.EventKey && e.Ch == 'q' {
          return "Everything went better than expected"
        }
        if e.Type == ui.EventResize {
          ui.Body.Width = ui.TermWidth()
          ui.Body.Align()
          go func() { redraw <- true }()
        }
      case <-done:
        return "Everything went better than expected"
      case e := <-error:
        return e
      case <-redraw:
        ui.Render(ui.Body)
        ui.Body.Align()
    }
  }
}

func main() {
  i_bet_its_not_ok := boot()
  log.Printf(i_bet_its_not_ok)
}