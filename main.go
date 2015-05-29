// Number One rule: it's compiling it's working, no tests needed
// Number Two: it's freaking imperative
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
  "errors"
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
  PossibleCurrentLanguageBecausePHPProgrammersDoesntKnowTheDifferenceBetweenBooleanAndStringThanks interface{} `json:"current_language,omitempty"`
  CurrentLanguage string
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

type LanguageByLevel []Language
func (v LanguageByLevel) Len() int { return len(v) }
func (v LanguageByLevel) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v LanguageByLevel) Less(i, j int) bool { 
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

func HandleMeLikeOneOfYourFrenchGirls(err error) {
  parUser := ui.NewPar(err.Error())
  parUser.Height = 3
  parUser.Border.Label = "Erro acontece nada ocorre feijoada"

  ui.Body.AddRows(
      ui.NewRow(ui.NewCol(12, 0, parUser)))
}

func main() {
  // yeah, if there's no username you'll see my profile :3
  var username string
  var numberOfLanguages int
  flag.StringVar(&username, "username", "schleumer", "your codeivate username")
  flag.IntVar(&numberOfLanguages, "len", 10, "number of languages to display")
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

  var startPoints float64
  firstRound := true

  update := func () {
    for {
      // restart body
      ui.Body = ui.NewGrid()
      ui.Body.Width = ui.TermWidth()

      client := &http.Client{}
      req, err := http.NewRequest("GET", fmt.Sprintf("http://codeivate.com/users/%s.json", username), nil)
      req.Close = true
      req.Header.Set("Content-Type", "application/json")
      req.Header.Set("User-Agent", "NOTICE ME SENPAI v0.1a")

      resp, err := client.Do(req)

      if err != nil {
        HandleMeLikeOneOfYourFrenchGirls(err)
        redraw <- true
        time.Sleep(time.Second * 5)
        continue
      }

      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        HandleMeLikeOneOfYourFrenchGirls(err)
        redraw <- true
        time.Sleep(time.Second * 5)
        continue
      }

      var statistic UserStatistic
      //err = json.Unmarshal([]byte(`
      //  {"name":"schleumer","level":"23.14","focus_level":"15.91","focus_points":2962,"max_streak":11094,"total_days_coded":380,"total_flow_states":981,"time_spent":1596843,"programming_now":false,"current_language":false,"streaking_now":false,"platforms":{"linux":{"percent_work":25.33,"points":810,"time":394999},"mac os x":{"percent_work":0.33,"points":10,"time":9571},"osx":{"percent_work":11.23,"points":359,"time":192277},"windows":{"percent_work":63.11,"points":2020,"time":999996}},"languages":{"ANY":{"level":"0.14","points":2},"AutoHotkey":{"level":"0.05","points":0},"Batch File":{"level":"0.41","points":6},"Build":{"level":"0.00","points":0},"C#":{"level":"0.05","points":0},"C++":{"level":"1.27","points":33},"CSS":{"level":"0.95","points":21},"Clojure":{"level":"0.55","points":9},"CoffeeScript":{"level":"2.73","points":112},"EJS":{"level":"1.14","points":27},"EJS_alternative":{"level":"0.18","points":2},"EJS_default":{"level":"0.23","points":3},"Elixir":{"level":"0.91","points":20},"Erlang":{"level":"0.05","points":0},"Go":{"level":"2.18","points":78},"GoSublime-Go":{"level":"1.95","points":63},"GoSublime-Template":{"level":"0.00","points":0},"Groovy":{"level":"0.00","points":0},"HOCON":{"level":"0.00","points":0},"HTML":{"level":"3.55","points":180},"HTML (Rails)":{"level":"0.36","points":6},"HTMLMustache":{"level":"0.05","points":1},"Haskell":{"level":"0.09","points":1},"JAVA":{"level":"0.41","points":6},"JSON":{"level":"1.41","points":38},"Jade":{"level":"2.68","points":109},"Java":{"level":"0.00","points":0},"JavaProperties":{"level":"0.00","points":0},"JavaScript":{"level":"5.27","points":364},"JavaScript (JSX)":{"level":"0.00","points":0},"JavaScriptNext":{"level":"0.00","points":0},"LESS":{"level":"2.45","points":93},"Lisp":{"level":"0.00","points":0},"LiveScript":{"level":"6.18","points":494},"Lua":{"level":"0.05","points":1},"Markdown":{"level":"2.05","points":68},"PHP":{"level":"5.73","points":424},"Plain":{"level":"0.32","points":5},"Plain Text":{"level":"2.50","points":97},"Play-Routes":{"level":"0.18","points":2},"Play-Scala-Template":{"level":"0.64","points":12},"Play2HtmlRouting":{"level":"0.00","points":0},"Play2Template":{"level":"0.05","points":0},"PowershellSyntax":{"level":"0.18","points":2},"Properties":{"level":"0.05","points":0},"Python":{"level":"1.55","points":44},"RegExp":{"level":"0.00","points":0},"Ruby":{"level":"1.64","points":49},"Ruby Slim":{"level":"0.14","points":2},"Rust":{"level":"0.55","points":10},"SQL":{"level":"0.14","points":1},"Sass":{"level":"0.05","points":0},"Scala":{"level":"7.18","points":651},"Shell-Unix-Generic":{"level":"0.32","points":5},"TypeScript":{"level":"0.09","points":1},"VimL":{"level":"0.00","points":0},"XML":{"level":"0.86","points":17},"YAML":{"level":"0.86","points":18},"configuration":{"level":"0.27","points":3},"laravel-blade":{"level":"2.45","points":93},"reStructuredText":{"level":"0.00","points":0},"yesod-hamlet":{"level":"0.00","points":0}}}
      //`), &statistic)
      err = json.Unmarshal(body, &statistic)
      if err != nil {
        //HandleMeLikeOneOfYourFrenchGirls(errors.New(fmt.Sprintf("Error on unmarshaling, probably the user %s doesn't exists [%s]", username, err.Error())))
        HandleMeLikeOneOfYourFrenchGirls(errors.New(fmt.Sprintf("[%s]", username, err.Error())))
        redraw <- true
        time.Sleep(time.Second * 60)
        continue
      }

      switch str := statistic.PossibleCurrentLanguageBecausePHPProgrammersDoesntKnowTheDifferenceBetweenBooleanAndStringThanks.(type) {
        case string:
          statistic.CurrentLanguage = str
        default:
          statistic.CurrentLanguage = "" 
      }

      userLevel, err := ParseLevel(statistic.Level)
      if err != nil {
        HandleMeLikeOneOfYourFrenchGirls(err)
        redraw <- true
        time.Sleep(time.Second * 5)
        continue
      }

      // i have no idea what i'm doing, but it's fucking hardcore
      hours := math.Floor(statistic.TimeSpent / 3600)
      minutes := math.Floor((statistic.TimeSpent - (hours * 3600)) / 60)
      var totalPoints float64
      var pointsSinceFirstRound float64

      for _, lang := range statistic.Languages {
        totalPoints += lang.Points
      }

      if firstRound {
        startPoints = totalPoints
      }

      pointsSinceFirstRound = totalPoints - startPoints

      parUser := ui.NewPar(
        fmt.Sprintf("Level: %d - Percent: %.0f\nTime: %.0f hours %.0f minutes - Current Language: %s\nPoints: %.0f(%.0f)", 
          userLevel.Number, 
          userLevel.Percent, 
          hours, 
          minutes, 
          statistic.CurrentLanguage, 
          totalPoints, 
          pointsSinceFirstRound))

      parUser.Height = 5
      parUser.Border.Label = "User Info"

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
      parWorkspace.Border.Label = "Workspace Info"

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

      sort.Sort(LanguageByLevel(languages))
      
      for _, lang := range languages[:numberOfLanguages] {
        g := ui.NewGauge()
        g.Percent = int(lang.Percent)
        g.Height = 3
        g.Border.Label = fmt.Sprintf("%s Level: %d - Points: %.0f", lang.Name, lang.Level, lang.Points)
        g.BarColor = ui.ColorGreen
  
        ui.Body.AddRows(
          ui.NewRow(ui.NewCol(12, 0, g)))
      }

      ui.Render(ui.Body)
      ui.Body.Align()

      redraw <- true

      // if everything went ok
      if firstRound {
        firstRound = false
      }
      time.Sleep(time.Second * 10)
    }
  }
  

  evt := ui.EventCh()

  ui.Render(ui.Body)
  go update()
  
  for {
    select {
      case e := <-evt:
        if e.Type == ui.EventKey && (e.Ch == 'q' || e.Ch == 'Q' /* HEHEHEHE */) {
          log.Print("Everything went better than expected")
          return
        }
        if e.Type == ui.EventResize {
          ui.Body.Width = ui.TermWidth()
          ui.Body.Align()
          go func() { redraw <- true }()
        }
      case <-done:
        log.Print("Everything went better than expected")
        return
      case e := <-error:
        log.Fatal(e)
        return
      case <-redraw:
        ui.Body.Align()
        ui.Render(ui.Body)
    }
  }
}