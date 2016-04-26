package main

type Config struct{}

func main() {
	/*svcConfig := &service.Config{
		Name:        "fakelive",
		DisplayName: "fakelive and opinion server",
		Description: "",
	}*/

	prg := &app{Quit: make(chan bool)}
	prg.run()
	select {}
	/*s, err := service.New(prg, svcConfig)
	       if err != nil {
		               log.Fatal(err)
		       }
	       if len(os.Args) > 1 {
		               err = service.Control(s, os.Args[1])
		               if err != nil {
			                       log.Fatal(err)
			               }
		               return
		       }

	       logger, err := s.Logger(nil)
	       if err != nil {
		               log.Fatal(err)
		       }
	       err = s.Run()
	       if err != nil {
		               logger.Error(err)
		       }*/

}
