package functions

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Struct for the long list format entries
type longListEntry struct {
	perm  string
	nlink string
	uid   string
	gid   string
	size  string
	time  string
	name  string
}

// Gets permission string
func permsString(info os.FileInfo) string {
	mode := info.Mode()
	// Make a string that has all set for comnparison
	all := "rwxrwxrwx"
	var buf [10]byte

	// Get first character
	if info.IsDir() {
		buf[0] = 'd'
	} else if mode.IsRegular() && mode&0111 != 0 {
		buf[0] = '-'
	} else {
		buf[0] = 'l'
	}

	// Get the remaining nine bits
	perms := mode & os.ModePerm

	// Check each bit by shifting a 1 to where we need to check then & to see if that bit is set
	for i, c := range all {
		if perms&(1<<uint(8-i)) != 0 {
			buf[i+1] = byte(c)
		} else {
			buf[i+1] = '-'
		}
	}
	return string(buf[:])
}

// Formats to human readable numbers when -h
func humanReadable(size int64) string {
	// List of prefixes to iterate through at each magnitude
	prefixes := []string{"", "K", "M", "G", "T", "P", "E", "Z", "Y", "R", "Q"}
	fsize := float64(size)

	// Keep dividing by 1024 to get to next magnitude
	i := 0
	for fsize >= 1024 && i < len(prefixes)-1 {
		fsize /= 1024
		i++
	}

	// Remove decimal if multiple digits
	if fsize >= 10 || i == 0 {
		return strconv.FormatInt(int64(math.Ceil(fsize)), 10) + prefixes[i]
	}

	// Keep one decimal place if not
	return strconv.FormatFloat(math.Ceil(fsize*10)/10, 'f', 1, 64) + prefixes[i]
}

// pads the string on the left with blank space (alligns right)
func padLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(" ", length-len(s))
	return padding + s
}

// pads the string on the right with blank space (alligns left)
func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(" ", length-len(s))
	return s + padding
}

// Get the number of blocks for the total blocks printed during -l
func getBlocks(dir string, entries []os.DirEntry) int64 {
	var blocks int64 = 0
	for _, entry := range entries {
		name := entry.Name()

		// Get file info and add blocks
		info, err := os.Lstat(filepath.Join(dir, name))
		if err == nil {
			stat := info.Sys().(*syscall.Stat_t)
			blocks += stat.Blocks
		}
	}
	return blocks / 2
}

// Long list format (-l/-n)
func llist(w io.Writer, dir string, names []string, flags []bool, blocks int64, dirs *[]string, useColor bool) error {
	// For max lengths of each category
	nlinksMax := 0
	uidsMax := 0
	gidsMax := 0
	sizesMax := 0

	var entries []longListEntry

	for _, name := range names {
		// Get file info
		info, err := os.Lstat(filepath.Join(dir, name))
		if err != nil {
			fmt.Fprint(os.Stderr, "gols: cannot access '"+name+"': No such file or directory\n")
			continue
		}
		mode := info.Mode()

		if info.IsDir() {
			// Note deeper directories to later iterate through if -R
			if flags[4] && name != "." && name != ".." {
				*dirs = append(*dirs, filepath.Join(dir, name))
			}
		}

		// Get stat
		stat := info.Sys().(*syscall.Stat_t)

		// Permissions
		perm := permsString(info)

		// Number of hard links
		nlink := strconv.FormatUint(uint64(stat.Nlink), 10)
		if len(nlink) > nlinksMax {
			nlinksMax = len(nlink)
		}

		// User ID
		uid := strconv.Itoa(int(stat.Uid))
		if !flags[2] {
			u, err := user.LookupId(uid)
			if err != nil {
				return err
			}
			uid = u.Username
		}
		if len(uid) > uidsMax {
			uidsMax = len(uid)
		}

		// Group ID
		gid := strconv.Itoa(int(stat.Gid))
		if !flags[2] {
			g, err := user.LookupGroupId(gid)
			if err != nil {
				return err
			}
			gid = g.Name
		}
		if len(gid) > gidsMax {
			gidsMax = len(gid)
		}

		// Sizes, convert to human readable if -h
		var size string
		if flags[3] {
			size = humanReadable(stat.Size)
		} else {
			size = strconv.FormatInt(stat.Size, 10)
		}
		if len(size) > sizesMax {
			sizesMax = len(size)
		}

		// Last modified time with different format for older than 180 days
		modified := time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
		var formatted string
		now := time.Now()

		// Get six months ago.
		// Apparently ls does not use 180 days for six months. I had some issues with files from 180 days before testing
		sixMonths := 15778476 * time.Second
		if modified.After(now.Add(-sixMonths)) {
			formatted = modified.Format("Jan _2 15:04")
		} else {
			formatted = modified.Format("Jan _2  2006")
		}

		// Get name in colors if necessary
		if useColor && info.IsDir() {
			// Print directories in blue
			name = Blue.ColorString(name)
		} else if useColor && mode.IsRegular() && mode&0111 != 0 {
			// Print files/executables in green
			name = Green.ColorString(name)
		} else if mode&os.ModeSymlink != 0 {
			// Handle links
			target, _ := os.Readlink(filepath.Join(dir, name))
			trgtInfo, err := os.Lstat(filepath.Join(filepath.Dir(name), target))
			if err != nil {
				// If can't read link, use red
				if useColor {
					name = Red.ColorString(name)
					target = Red.ColorString(target)
				}
			} else {
				if useColor {
					// Color depending on info
					trgtMode := trgtInfo.Mode()
					name = Cyan.ColorString(name)

					if trgtInfo.IsDir() {
						target = Blue.ColorString(target)
					} else if trgtMode.IsRegular() && trgtMode&0111 != 0 {
						target = Green.ColorString(target)
					} else if trgtMode&os.ModeSymlink != 0 {
						target = Cyan.ColorString(target)
					}
				}
			}
			// Format link
			name = name + " -> " + target
		}
		// Add to entries because we need to write after we have determined column width
		entries = append(entries, longListEntry{perm: perm, nlink: nlink, uid: uid, gid: gid, size: size, time: formatted, name: name})
	}

	// Include total blocks if going in directory. Reformat if -h
	if blocks >= 0 {
		if flags[3] {
			Print(w, "total "+humanReadable(blocks*1024)+"\n")
		} else {
			Print(w, "total "+strconv.FormatInt(blocks, 10)+"\n")
		}
	}

	// Print all entries with the correct padding
	for _, entry := range entries {
		str := entry.perm + " " + padLeft(entry.nlink, nlinksMax) + " " +
			padRight(entry.uid, uidsMax) + " " + padRight(entry.gid, gidsMax) + " " + padLeft(entry.size, sizesMax) +
			" " + entry.time + " " + entry.name + "\n"
		Print(w, str)
	}
	return nil
}

// Performs ls on a directory
func lsDir(dir string, w io.Writer, flags []bool, useColor bool) {
	if flags[4] {
		// Print dir name if -R
		if dir != "." && strings.Split(dir, string(filepath.Separator))[0] != ".." {
			Print(w, "."+string(filepath.Separator)+dir+":\n")
		} else if dir == "." {
			Print(w, ".:\n")
		} else {
			Print(w, dir+":\n")
		}

	}

	// Get files in directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Get the number of blocks before filtering if -l or -n
	var blocks int64
	if flags[1] || flags[2] {
		blocks = getBlocks(dir, entries)
	}

	names := []string{}
	if flags[0] {
		// Include parent and working directories if -a
		names = []string{".", ".."}
	} else {
		// Filter entries if not -a
		entries = dirFilter(entries)
	}

	// Get names
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	// Explcitiely sort names (not sure if this is required or not by the directions)
	sort.Strings(names)

	// Create directories to iterate through for -R
	dirs := []string{}
	if flags[1] || flags[2] {
		// Call llist if -l/-n
		err := llist(w, dir, names, flags, blocks, &dirs, useColor)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	} else {
		// Otherwise:
		for _, name := range names {
			// Get file info
			info, err := os.Lstat(filepath.Join(dir, name))
			if err != nil {
				fmt.Fprintln(os.Stderr, "gols: cannot access '"+name+"': No such file or directory")
				continue
			}
			mode := info.Mode()

			if info.IsDir() {
				// Note deeper directories to later iterate through if -R
				if flags[4] && name != "." && name != ".." {
					dirs = append(dirs, filepath.Join(dir, name))
				}
			}

			if useColor && info.IsDir() {
				// Print directories in blue
				Blue.ColorPrint(w, name+"\n")
			} else if useColor && mode.IsRegular() && mode&0111 != 0 {
				// Print files/executables in green
				Green.ColorPrint(w, name+"\n")
			} else if mode&os.ModeSymlink != 0 {
				// Print links in cyan
				Cyan.ColorPrint(w, name+"\n")
			} else {
				// Print anything else in default color
				Print(w, name+"\n")
			}
		}
	}

	if flags[4] {
		// If -R recursively process each directory
		for _, d := range dirs {
			// Add a new line for readability
			Print(w, "\n")

			lsDir(d, w, flags, useColor)
		}
	}
}

// ls with flags implementation
func LS(w io.Writer, args []string, flags []bool, useColor bool) error {
	if len(args) == 0 {
		// Perform ls on current directory if no arguments are given
		lsDir(".", w, flags, useColor)
	} else {
		// Split args into files and directories
		var files []string
		var dirs []string
		for _, arg := range args {
			info, err := os.Lstat(arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, "gols: cannot access '"+arg+"': No such file or directory")
				continue
			}

			if info.IsDir() {
				dirs = append(dirs, arg)
			} else {
				files = append(files, arg)
			}
		}

		// Explicitely sort all sets of data
		sort.Strings(dirs)
		sort.Strings(files)

		if flags[1] || flags[2] {
			err := llist(w, ".", files, flags, -1, &[]string{}, useColor)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
			}
		} else {
			// Perform ls on files
			for _, file := range files {
				if useColor {
					// Get file info
					info, err := os.Lstat(file)
					if err != nil {
						fmt.Fprintln(os.Stderr, "gols: cannot access '"+file+"': No such file or directory")
						continue
					}
					mode := info.Mode()

					if mode.IsRegular() && mode&0111 != 0 {
						// Print files/executables in green
						Green.ColorPrint(w, file+"\n")
					} else if mode&os.ModeSymlink != 0 {
						// Print links in cyan
						Cyan.ColorPrint(w, file+"\n")
					} else {
						// Print anything else in default color
						Print(w, file+"\n")
					}
				} else {
					// Print in default color if not using colors
					Print(w, file+"\n")
				}
			}
		}

		// Perform ls on each directory
		for _, dir := range dirs {
			if len(args) > 1 && !flags[4] {
				// Add a header in default color if there were multiple arguments
				/*
					Please note that while the assignement asks for this if multiple
					directories are targets, Unix ls just does this if there are
					multiple targets. That is how it is implemented here.
				*/
				Print(w, "\n"+dir+": \n")
			}
			lsDir(dir, w, flags, useColor)
		}
	}
	return nil
}
