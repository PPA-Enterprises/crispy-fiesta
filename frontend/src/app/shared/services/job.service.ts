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


@Injectable({
  providedIn: 'root'
})
export class JobService {

  jobs: Job[] = [
    { 
      id: 0,
      client_name: "Tristan Hull",
      client_phone: "661-208-1140",
      date: "December 12th, 2020", 
      car_info: "2003 Saturn Vue",
      appointment_info: "The Appointment is at 12/27/2020 5:00pm",
      notes: "These are some example notes",
      tag: "OPEN"
    },
    { 
      id: 1,
      client_name: "Tristan Hull",
      client_phone: "661-208-1140",
      date: "December 12th, 2020", 
      car_info: "2003 Saturn Vue",
      appointment_info: "The Appointment is at 12/27/2020 5:00pm",
      notes: "These are some example notes",
      tag: "CLOSED"
    },
    { 
      id: 2,
      client_name: "Tristan Hull",
      client_phone: "661-208-1140",
      date: "December 12th, 2020", 
      car_info: "2003 Saturn Vue",
      appointment_info: "The Appointment is at 12/27/2020 5:00pm",
      notes: "These are some example notes",
      tag: "NEW"
    },
  ];

  constructor(private http: HttpClient) { }

  createJob (job: Job): Job {
    // return this.http.post<Job>("http://192.168.1.21:8080/api/v1/jobs", job)
    this.jobs.push(job);
    return job;
  }

  getAllJobs(): Job[]{
    return this.jobs;
  }

  getJobById(id: number): Job {
    return this.jobs.find(job => job.id === id);
  }

  editJobById(id: number, job: Job): Job {
    this.jobs[this.jobs.findIndex(job => job.id === id)] = job;
    return job;
  }

  deleteJobById(id: number): number {
    this.jobs.splice(this.jobs.findIndex(job => job.id === id), 1);
    return id;
  }
}
