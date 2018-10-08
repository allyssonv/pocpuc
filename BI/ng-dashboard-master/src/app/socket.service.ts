import { Injectable } from '@angular/core'
import { Observable } from 'rxjs'
import * as io from 'socket.io-client'

@Injectable()
export class SocketService {

  private SERVER_URL = 'http://localhost:3000'
  private socket

  
  getSales(){

    let observable = new Observable( observer => {
      this.socket = io(this.SERVER_URL)
      this.socket.on('sales', (data) => {
        observer.next(data)
      })
      return () => {
        this.socket.disconnect()
      }
    })
    return observable

  }

  

}
