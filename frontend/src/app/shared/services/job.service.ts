import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { startOfDay, endOfDay, subDays, addDays, endOfMonth, isSameDay, isSameMonth, addHours } from 'date-fns';
import { Job } from '../models/job.model';

const colors: any = {
  red: {
    primary: '#ad2121',
    secondary: '#FAE3E3',
  },
  blue: {
    primary: '#1e90ff',
    secondary: '#D1E8FF',
  },
  yellow: {
    primary: '#e3bc08',
    secondary: '#FDF1BA',
  },
};


@Injectable({
  providedIn: 'root'
})
export class JobService {

  jobs: Job[] = [
    { 
      id: '0',
      client_name: "Tristan Hull",
      client_phone: "661-208-1140",
      car_info: "2003 Saturn Vue",
      appointment_info: "The Appointment is at 12/27/2020 5:00pm",
      notes: "These are some example notes",
      tag: "OPEN",
      start: subDays(startOfDay(new Date()), 1),
      end: addDays(new Date(), 1),
      color: colors.red,
      title: '2003 Saturn Vue - Tristan Hull',
    },
    { 
      id: '1',
      client_name: "Frank Hull",
      client_phone: "661-208-1140",
      car_info: "2003 Saturn L200",
      appointment_info: "The Appointment is at 12/27/2020 5:00pm",
      notes: "waddup",
      tag: "CLOSED",
      start: subDays(startOfDay(new Date()), 2),
      end: addDays(new Date(), 2),
      color: colors.yellow,
      title: '2003 Saturn L200 - Frank Hull',
    },
    { 
      id: '2',
      client_name: "Tristan Hull",
      client_phone: "661-208-1140",
      car_info: "2020 Audi S5",
      appointment_info: "The Appointment is at 12/27/2020 5:00pm",
      notes: "These are some example notes",
      tag: "NEW",
      start: subDays(startOfDay(new Date()), 5),
      end: addDays(new Date(), 5),
      color: colors.yellow,
      title: '2020 Audi S5 - Frank Swartz',
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
    return this.jobs.find(job => Number(job.id) === id);
  }

  editJobById(id: number, job: Job): Job {
    this.jobs[this.jobs.findIndex(job => Number(job.id) === id)] = job;
    return job;
  }

  deleteJobById(id: number): number {
    this.jobs.splice(this.jobs.findIndex(job => Number(job.id) === id), 1);
    return id;
  }
}
