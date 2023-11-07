package main

import (
	"github.com/gobwas/ws"
	"log"
	"net"
	"net/http"
)

func main() {
	s := newServer()

	err := http.ListenAndServe(":8585", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			panic(err)
		}
		//defer conn.Close()
		go func() {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					log.Println("main->defer : ", err)
				}
			}(conn)

			for {
				s.run()

			}

		}()
		// s.run()

		/*c := s.NewClient(conn)
		go c.ReadInput()*/
		go s.NewClient(conn).ReadInput()

	}))
	if err != nil {
		log.Println("main->server : ", err)
	}
}
