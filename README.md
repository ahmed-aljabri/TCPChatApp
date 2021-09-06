#TCPChatApp

    - A chat application program that allows multiple clients to create chat rooms and exchange messages in a bi-directional manner.

    - Note (03.09.2021):
        This application is functional but yet under development, the curent state of the app and the future updated features will be listed below.

        USAGE:
            - Navigate to the Server directory and run: "go build ."
            - run the executable and specify the listening port: "./Server.exe 8888"
            - Now navigate to the client directory and run: go build .
            - run the executable and specify the IP/hostname & port: "localhost:8888"

        AJ

    #V1.0 [03.09.2021]:
        - Bi-directional chat between clients on individually created chat rooms.
        - Data transfer is on raw TCP

    #2.0 [06.09.2021]:
        - Added support for TLS option in addition to raw TCP
        - fixed some issues with message broadcasting
    

#
