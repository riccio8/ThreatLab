/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */



#include <iostream>
#include <boost/asio.hpp>
#include <cstring>

using namespace std;
using namespace boost::asio;
using ip::tcp;

class tcp_sock {
public:
    tcp_sock(const string& ip, const string& port, const string& msg) {
        while(true){
         io_context io;

          tcp::resolver resolver(io);
          tcp::resolver::results_type endpoints = resolver.resolve(ip, port);


          tcp::socket socket(io);
          connect(socket, endpoints);

          boost::system::error_code error;
          write(socket, buffer(msg), error);

          if (!error) {
              cout << "Message sent successfully!" << endl;
          } else {
              cout << "Error occurred: " << error.message() << endl;
          }
        }

    }

    ~tcp_sock() {
    }
};

int main() {
    string ip = "127.0.0.1";  
    string port = "8080";     
    string msg = "Hello from client!";  

    tcp_sock Tcp_con(ip, port, msg);

    return 0;
}
