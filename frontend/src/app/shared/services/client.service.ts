import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { from, Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { Job } from '../models/job.model';
import { Client } from '../models/client.model';
import { JobService } from '../services/job.service';

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
            jobs: this.jobService.getAllJobs()
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
    constructor(private http: HttpClient, private jobService: JobService) { }

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
