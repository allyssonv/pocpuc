/**
 * Created by andrew.yang on 5/18/2017.
 */

import {OnInit, Component} from "@angular/core";
import { SocketService } from '../../socket.service';

@Component({
    selector: 'home',
    templateUrl: './home.component.html',
    providers: [SocketService]
})
export class HomeComponent implements OnInit {

    price:number
    total:number = 0
    info:any = {}
    comments:string[] = []
    items:string[] = []
    status: string

    data:any = {}

    comm:any = {}

    shippingStatus:any[] = ["shipped", 'processing', 'delivered', 'in transit']

    commenter:any[] = ["Sandye Willats", "Matt McGeechan",
                       "Florencia Mingaud", "Omar Stilldale",
                       "Merwin Codeman", "Joe Higashi",
                       "Terry Bogard", "Hanzo Hasashi"]


    constructor(private socket: SocketService) {
        setInterval( ()=> {
              this.items.length = 0
              this.comments.length = 0  
        }, 60000)

     }

    ngOnInit() {

       this.socket.getSales().subscribe( message => {

       this.info = JSON.parse(message.toString())
       //========================================================
       this.total = this.total + parseInt(this.info.Price) // total das vendas
       //========================================================
       this.comm.name = this.commenter[Math.floor(Math.random() * this.commenter.length)]
       this.comm.comment =  this.info.Comment 
       this.comm.dateComment = new Date().toLocaleString("pt-BR")
       this.comments.push( this.comm )
       this.comm = {}
       //========================================================
       this.data.order = this.info.OrderNo
       this.data.name = this.info.Name
       this.data.date = new Date().toLocaleString("pt-BR")
       this.data.status = this.shippingStatus[Math.floor(Math.random() * this.shippingStatus.length)]
       this.items.push(this.data)
       this.data = {}
       //========================================================
    })



    }
}