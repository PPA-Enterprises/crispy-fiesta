import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

export class Job {
  public client_name: string
  public client_phone: string;
  public car_info: string;
  public appointment_info: string;
  public notes: string;
}


@Injectable({
  providedIn: 'root'
})
export class JobService {

  constructor(private http: HttpClient) { }

  createJob (job: Job): Observable<Job> {
    return this.http.post<Job>("http://192.168.1.21:8080/api/v1/jobs", job)
  }

}
