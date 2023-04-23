package cmd

func Run() {
	/*fmt.Println("服务启动中···")
	listener, err := net.Listen("tcp", "0.0.0.0:6060")
	if err != nil {
		print("服务启动失败")
		return
	}
	fmt.Println(listener)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	for {
		fmt.Println("服务启动中···")
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		fmt.Println(conn)
	}*/
}

func init() {
	getPwd()
}
