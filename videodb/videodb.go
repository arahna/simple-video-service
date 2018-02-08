package videodb

type Video struct {
	Id        string
	Name      string
	Duration  int
}

var videoList = [3]Video{
	Video{
		"d290f1ee-6c54-4b01-90e6-d701748f0851",
		"Black Retrospetive Woman",
		15,
	},
	Video{
		"hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345",
		"Go Rally TEASER-HD",
		41,
	},
	Video{
		"sldjfl34-dfgj-523k-jk34-5jk3j45klj34",
		"Танцор",
		92,
	},
}

func GetAll() []Video {
	return videoList[:]
}

func Find(id string) (Video, bool) {
	for _, video := range videoList {
		if video.Id == id {
			return video, true
		}
	}
	return Video{}, false
}