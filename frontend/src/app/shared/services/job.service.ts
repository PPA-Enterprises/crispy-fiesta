import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { startOfDay, endOfDay, subDays, addDays, endOfMonth, isSameDay, isSameMonth, addHours } from 'date-fns';
import { Job } from '../models/job.model';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class JobService {
	url ="http://ppaenterprises.com:8888/"

  constructor(private http: HttpClient) { }

  createJob (job: Job): Observable<Job> {
    return this.http.post<any>(this.url + "api/v1/jobs/", job)
    .pipe(map((response) => {
			if(response.success) {
				return response.payload as Job;
			} else {
				return null;
			}
		  }));
  }

  getAllJobs(): Observable<Job[]> {
    return this.http.get<any>(this.url + "api/v1/jobs/")
    .pipe(map((response) => {
			if(response.success) {
				return response.payload as Job[];
			} else {
				return [];
			}
		  }));
  }

  getJobById(id: string): Observable<Job> {
    return this.http.get<any>(this.url + "api/v1/jobs/id/"+id)
    .pipe(map((response) => {
			if(response.success) {
        response.payload.start = new Date(response.payload.start);
        response.payload.end = new Date(response.payload.end);
				return response.payload as Job;
			} else {
				return null;
			}
		  }));
  }

  editJobById(id: string, job: Job): Observable<Job> {
    return this.http.patch<any>(this.url + "api/v1/jobs/"+id, job)
    .pipe(map((response) => {
			if(response.success) {
				return response.payload as Job;
			} else {
				return null;
			}
		  }));
  }

  deleteJobById(id: string): Observable<any> {
    return this.http.delete<any>(this.url + "api/v1/jobs/"+id)
    .pipe(map((response) => {
			if(response.success) {
				return true;
			} else {
				return null;
			}
		}));
  }
}
