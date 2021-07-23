module github.com/Postcord/router

go 1.16

require (
	github.com/Postcord/interactions v0.0.13
	github.com/Postcord/objects v0.0.11-0.20210606154159-a696ad5ad3cd
	github.com/Postcord/rest v0.0.5-0.20210607004326-f5827f030be6
)

replace (
	github.com/Postcord/objects v0.0.11-0.20210606154159-a696ad5ad3cd => git.kelwing.dev/Postcord/objects v0.0.12-0.20210723081135-f6c9b9cad7a5
	github.com/Postcord/rest v0.0.5-0.20210607004326-f5827f030be6 => git.kelwing.dev/Postcord/rest v0.0.6-0.20210723081922-d28de31375a1
)
