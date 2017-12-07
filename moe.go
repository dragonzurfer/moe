package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jaytaylor/html2text"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// command params
var (
	name, AnimeURL, AnimeVideoURL, seasonal, video     string
	score, rank, synopsis, info, songs, EP, aired, all bool
	MALsearch                                          = "https://myanimelist.net/search/all?q="
	VIDEOsearch                                        = "https://9anime.is/search?keyword="
)

// results
var (
	scoreres, rankres, statres, OPres, EDres, EPres, airedres string
	synopsisres, songsres, seasonalres                        []string
	infores                                                   = make(map[string]string)
)

//colors
var (
	green         = color.New(color.FgHiGreen)
	boldcyan      = color.New(color.FgCyan, color.Bold)
	boldred       = color.New(color.FgRed, color.Bold)
	boldblue      = color.New(color.FgHiBlue, color.Bold)
	boldyellow    = color.New(color.FgYellow, color.Bold)
	boldwhite     = color.New(color.FgHiWhite, color.Bold)
	boldgreen     = color.New(color.FgHiGreen, color.Bold)
	italicmagenta = color.New(color.FgHiMagenta)
	italicblue    = color.New(color.FgBlue, color.Italic)
)

// bind flags to params
func bindFlags() {
	flag.StringVar(&name, "name", "", "Give Name ex: DeathNode, \"Your Lie In April\"")
	flag.StringVar(&seasonal, "seasonal", "", "<SEASON> <YEAR> ex: summer 2017, winter 2016 or Just leave blank for current season")
	flag.StringVar(&video, "video", "", "<EPISODE NUMBER> ex: 1, 9 etc or \"all\" to get all the episodes")
	flag.BoolVar(&score, "score", false, "Get Score")
	flag.BoolVar(&rank, "rank", false, "Get Rank")
	flag.BoolVar(&synopsis, "synopsis", false, "Get Synopsis")
	flag.BoolVar(&info, "info", false, "Get information")
	flag.BoolVar(&songs, "songs", false, "Get all the Opening and Ending song names")
	flag.BoolVar(&EP, "EP", false, "Get number of episodes")
	flag.BoolVar(&aired, "aired", false, "Get the aired date")
	flag.BoolVar(&all, "all", false, "Get All Params")
	flag.Parse()
}

// Get HTML page as string
func getContent(URL string) (string, bool) {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error fetching page")
		return "", true
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ret := string(body)
	if err != nil {
		fmt.Println("Error :( Try Again")
		return "", true
	}
	return ret, false
}

// check if ':' exists
func check(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return true
		}
	}
	return false
}

// replace
func Rep(s *string, rep [][]string) {
	var temp string = *s
	for i := 0; i < len(rep); i++ {
		temp = strings.Replace(temp, rep[i][0], rep[i][1], -1)
	}
	*s = temp
}

// Print to terminal
func PrintParams() {

	if info || all {
		boldblue.Printf("Information\n------------\n")
		for key, value := range infores {
			boldcyan.Printf("%v", key)
			for i := 0; i < 9-len(key); i++ {
				fmt.Printf(" ")
			}
			if key == "Score" {
				boldgreen.Printf(":  %v\n", scoreres)
				continue
			}
			boldwhite.Printf(": %v\n", value)
		}
		fmt.Printf("\n")
	}

	if seasonal != "" {
		boldblue.Printf("Animes of %v season \n-------------------\n", seasonal)
		for index, anime := range seasonalres {
			boldwhite.Printf("%v.", index)
			fmt.Printf(" %v\n", anime)
		}
		fmt.Printf("\n")
	}

	if synopsis || all {
		boldblue.Printf("Synopsis\n------------\n")
		for i := 0; i < len(synopsisres); i++ {
			boldwhite.Printf(synopsisres[i])
		}
		fmt.Printf("\n\n")
	}

	if songs || all {
		boldblue.Printf("Songs (OP's & EP's)\n----------------\n")
		for _, song := range songsres {
			italicmagenta.Println(song)
		}
		fmt.Printf("\n")
	}

	if score || all {
		boldblue.Printf("Score\n------------\n")
		boldgreen.Println(scoreres)
		fmt.Printf("\n")
	}

	if rank || all {
		boldblue.Printf("Ranked\n------------\n")
		boldred.Println(rankres)
		fmt.Printf("\n")
	}

	if aired || all {
		boldblue.Printf("Aired\n------------\n")
		boldwhite.Println(airedres)
		fmt.Printf("\n")
	}

	if EP || all {
		boldblue.Printf("Episodes\n------------\n")
		boldwhite.Println(EPres)
		fmt.Printf("\n")
	}
}

// Check if No results found
func emptyResult(length int, FLAG, query string) bool {
	if length == 0 {
		boldred.Printf("Could not find any results for %v: %v\n", FLAG, query)
		return true
	}
	return false
}

func downloadVideo() {
	// Get number of episodes on AnimeVideoURL page

	resp, err := getContent(AnimeVideoURL)

	if err {
		return
	}

	htmlcontent, errstr := html2text.FromString(resp)

	if errstr != nil {
		panic(errstr)
	}
	htmlcontent = strings.Replace(htmlcontent, " ", "", -1)
	regexServerG4 := (`ServerG4[^\*]*|\*(\d+)\(([^\)]*)`)
	serverg4re := regexp.MustCompile(regexServerG4)

	results := serverg4re.FindAllStringSubmatch(htmlcontent, -1)
	episodeURL := make(map[int]string)
	videoSite := `https://9anime.is`
	totalEpisodes := -1

	// map episode number to a url for download
	for _, result := range results {
		if result[1] == "" {
			continue
		}

		episodeNumber, converror := strconv.Atoi(result[1])

		if converror != nil {
			fmt.Println("Error downloading ...")
			return
		}

		if episodeURL[episodeNumber] != "" {
			break
		}

		episodeURL[episodeNumber] = videoSite + result[2]
		totalEpisodes = episodeNumber
	}

	if totalEpisodes == -1 {
		fmt.Println("Error downloading Anime could'nt find any downloads")
		return
	}

	boldgreen.Println(totalEpisodes)
	for k, v := range episodeURL {
		fmt.Println(k, v)
	}
}

// fetch seasonl animes
func fetchDetailsSeason(seasonal string) bool {
	regexseason := `TV \(New\)(.|\n)*?ONA`
	seasonalre := regexp.MustCompile(regexseason)
	var seasonURL string
	if seasonal == "CURRENT" {
		seasonURL = "https://myanimelist.net/anime/season"
	} else {
		temp := strings.Split(seasonal, " ")
		seasonURL = "https://myanimelist.net/anime/season/" + temp[1] + "/" + temp[0]
	}

	resp, err := getContent(seasonURL)
	if err {
		return false
	}

	cleaninfo, errstr := html2text.FromString(resp, html2text.Options{PrettyTables: true})
	if errstr != nil {
		panic(err)
	}

	cleaninfo = seasonalre.FindString(cleaninfo)
	regex := `https://myanimelist.net/anime/[0-9]*/([^\s]*)`
	re := regexp.MustCompile(regex)
	animes := re.FindAllString(cleaninfo, -1)
	rep := [][]string{{"-", " "}, {"_", " "}}
	m := make(map[string]bool)
	for _, anime := range animes {
		res := strings.Split(anime, "/")
		Rep(&res[5], rep)
		if !m[res[5]] {
			seasonalres = append(seasonalres, res[5])
			m[res[5]] = true
		}
	}

	return true
}

// Get anime details from MAL
func fetchDetails() bool {

	resp, err := getContent(AnimeURL)
	if err {
		return false
	}
	extractregex := regexp.MustCompile(">(.|\n)*?<")

	// Extract Synopsis
	regexsynopsis := `<span\sitemprop="description">(.|\n)*?</span>`
	regexscore := `[0-9]\.[0-9]{2,}`
	regexinfo := `<h2[^>]*>(.|\n)*?</h2>(.|\n)*?(<div[^>]*>(.|\n)*?</div>\s*)+`
	regexop := `<span\s+class="theme-song">(.|\n)*?</span>`
	regexrank := `Ranked[^<]*<strong>(#[\d]+)</strong>`
	synopsisre := regexp.MustCompile(regexsynopsis)
	infore := regexp.MustCompile(regexinfo)
	scorere := regexp.MustCompile(regexscore)
	opre := regexp.MustCompile(regexop)
	rankre := regexp.MustCompile(regexrank)

	// Extract synopsis
	result := synopsisre.FindAllString(resp, 1)
	synopsisres = extractregex.FindAllString(result[0], -1)

	if emptyResult(len(synopsisres), "-synopsis", "true") {
		return false
	}

	synopsisres[0] = synopsisres[0][1 : len(synopsisres[0])-1]
	rep := [][]string{{"<", ""}, {">", ""}}

	for i := 0; i < len(synopsisres); i++ {
		Rep(&synopsisres[i], rep)
		synopsisres[i] = html.UnescapeString(synopsisres[i])
	}

	// Extract score
	result = scorere.FindAllString(resp, 1)

	if emptyResult(len(result), "-score", "true") {
		return false
	}

	scoreres = result[0]

	// Extract info
	dirtyinfo := strings.Join(infore.FindAllString(resp, 3), "")
	cleaninfo, errstr := html2text.FromString(dirtyinfo)
	if errstr != nil {
		panic(err)
	}
	splitcleaninfo := strings.Split(cleaninfo, "\n")
	m := make(map[string]string)

	for _, str := range splitcleaninfo {
		checkres := check(str)
		if checkres {
			r := strings.Split(str, ":")
			m[r[0]] = r[1]
		}
	}

	cleaner := `((\([\d\w/\s_\-]*\))|(\(\s+[\w]+$))`
	cleanere := regexp.MustCompile(cleaner)
	for k, _ := range m {
		infores[k] = cleanere.ReplaceAllString(m[k], "")
	}

	// Extract OP's and ED's
	songsres = opre.FindAllString(resp, -1)

	for i := 0; i < len(songsres); i++ {
		songsres[i] = html.UnescapeString(extractregex.FindString(songsres[i]))
		Rep(&songsres[i], rep)
	}

	// Extract rank
	rankres = extractregex.FindString(html.UnescapeString(rankre.FindString(resp)))
	Rep(&rankres, rep)

	// Extract number of episodes
	EPres = m["Episodes"]

	// Extract aired date
	airedres = m["Aired"]

	// Extract seasonal anime
	if seasonal == "CURRENT" {
		seasonal = ""
	}

	return true
}

// Search given a name
// true value indicates matching name found
// shows search results if not found
func Search() bool {
	searchURL := MALsearch + name
	name = strings.ToLower(name)

	//make GET request
	resp, err := getContent(searchURL)
	if err {
		return false
	}

	regex := `<article>(.|\n)*?</article>`
	re := regexp.MustCompile(regex)
	results := re.FindString(resp)

	regex = `https://myanimelist.net/anime/[0-9]*/([^"/]*)`
	re2 := regexp.MustCompile(regex)
	results2 := re2.FindAllStringSubmatch(results, -1)
	var foundAnime, foundVideoAnime bool

	animeUrlMap := make(map[string]bool)
	for _, res := range results2 {
		res[1] = strings.Replace(res[1], "_", " ", -1)
		res[1] = strings.Replace(res[1], "  ", " ", -1) // replace double space
		res[1] = strings.ToLower(res[1])
		if res[1] == name {
			// set anime url to fetch results
			AnimeURL = res[0]
			foundAnime = true
			break
		}
		animeUrlMap[res[1]] = true
	}

	if !foundAnime && video == "" {
		index := 0
		boldyellow.Println("Did you mean :")
		boldyellow.Println("---------------")
		for key := range animeUrlMap {
			index += 1
			green.Printf("%v.", index)
			fmt.Printf(" %v\n", key)
		}
	}

	// if no video parameters set return true if AnimeURL was found
	if video == "" {
		if foundAnime {
			return true
		}
		return false
	}

	// Search for episode videos
	searchName := strings.Replace(name, " ", "%20", -1)
	searchURL = VIDEOsearch + searchName
	resp, err = getContent(searchURL)

	if err {
		return false
	}

	regex = `https://9anime.is/watch/([^"]*)`
	re3 := regexp.MustCompile(regex)
	resultsVideo := re3.FindAllStringSubmatch(resp, -1)

	if emptyResult(len(resultsVideo), "-video", name) {
		return false
	}

	searchName = strings.Replace(searchName, "%20", "-", -1)
	animeVideoUrlMap := make(map[string]bool)
	for _, anime := range resultsVideo {
		dotpos := 0
		if len(anime) < 2 {
			continue
		}
		cleanname := ""
		for i := 0; i < len(anime[1]); i++ {
			if anime[1][i] == '.' {
				dotpos = i
				cleanname = strings.Replace(anime[1][:dotpos], "-", " ", -1)
				animeVideoUrlMap[cleanname] = true
				break
			}
		}
		if cleanname == searchName {
			foundVideoAnime = true
			AnimeVideoURL = anime[0]
			return true
		}
	}

	if !foundVideoAnime {
		index := 0
		boldyellow.Println("Video results found for :")
		boldyellow.Println("---------------")
		for key := range animeVideoUrlMap {
			index += 1
			green.Printf("%v.", index)
			fmt.Printf(" %v\n", key)
		}
	}
	return false
}

func main() {
	bindFlags()
	if name != "" {
		if !synopsis && !score && !rank && !info && !EP && !aired && !songs && !all && len(video) == 0 {
			boldred.Println("No params found")
			return
		}
		success := Search()
		if success {
			if fetchDetails() {
				PrintParams()
			}
			if video != "" {
				downloadVideo()
			}
		}
	} else if seasonal != "" {
		if fetchDetailsSeason(seasonal) {
			PrintParams()
		}
	}
}
