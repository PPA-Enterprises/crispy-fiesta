import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { from, Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { Job } from '../models/job.model';
import { Client } from '../models/client.model';
import { JobService } from '../services/job.service';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ClientService {
    constructor(private http: HttpClient, private jobService: JobService) { }

    getAllClients(): Observable<Client[]> {
      return this.http.get<any>("http://ppaenterprises.com:8888/api/v1/clients")
      .pipe(map((response) => {
        if(response.success) {
          return response.payload as Client[];
        } else {
          return [];
        }
        }));
    }

    getClientById(id: string): Observable<Client> {
      return this.http.get<any>("http://ppaenterprises.com:8888/api/v1/clients/id/"+id)
      .pipe(map((response) => {
        if(response.success) {
          return response.payload as Client;
        } else {
          return null;
        }
        }));
    }

    editClientById(id: string, client: Client): Observable<Client> {
      return this.http.patch<any>("http://ppaenterprises.com:8888/api/v1/clients/"+id, client)
      .pipe(map((response) => {
        if(response.success) {
          return response.payload as Client;
        } else {
          return null;
        }
        }));
      }

}
