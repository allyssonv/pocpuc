import { Component, OnInit } from '@angular/core';
import {Login} from "./models/login";
import { SocketService } from './socket.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [SocketService]
})
export class AppComponent implements OnInit{

  info:any = {}

  foo:string 

  loginInfo:Login = {
      first_name:'Mr.',
      last_name:'Anderson',
      avatar: 'admin.png',
      title:'Senior Developer'
  };

  constructor(private socket: SocketService){}

  ngOnInit() {

      this.socket.getSales().subscribe( message => {

       //09090hjingf
      
    })



  }

}
