import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

export class Job {
    public id: number;
    public client_name: string
    public client_phone: string;
    public car_info: string;
    public appointment_info: string;
    public notes: string;
    public tag: string;
    public date: string;
  }

export class Client {
    public id: number;
    public name: string
    public phone: string;
    public email: string;
    public jobs: Job[];
}

@Injectable({
  providedIn: 'root'
})
export class ClientService {
    clients: Client[] = [
        {
            id: 0, 
            name: "Tristan Hull", 
            email: "tristan@kenpokicks.com", 
            phone: "6612081140", 
            jobs: []
        }, 
        {
            id: 1, 
            name: "Neena Romero", 
            email: "neena@kenpokicks.com", 
            phone: "6612081140", 
            jobs: []
        },
        {
            id: 2, 
            name: "Joshua Benz", 
            email: "joshua@kenpokicks.com", 
            phone: "6612081140", 
            jobs: []
        },
    ];
    constructor(private http: HttpClient) { }

    getAllClients(): Client[] {
        return this.clients;
    }

    getClientById(id: number): Client {
        return this.clients.find(client => client.id === id);
    }

    editClientById(id: number, client: Client): Client {
        this.clients[this.clients.findIndex(client => client.id === id)] = client;
        return client;
      }

}
