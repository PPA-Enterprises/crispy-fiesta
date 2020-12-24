import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { JobService } from '../shared/services/job.service'

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

@Component({
  selector: 'app-page',
  templateUrl: './jobs.component.html',
  styleUrls: ['./jobs.component.scss']
})

export class JobsComponent implements OnInit {
  newJobs: Job[] = [];
  openJobs: Job[] = [];
  closedJobs: Job[] = [];
  allJobs: Job[] = [];
  constructor(private jobService: JobService) {
    this.allJobs = this.jobService.getAllJobs();

    for(let job of this.allJobs) {
      if(job.tag == "NEW") {
        this.newJobs.push(job);
      } else if(job.tag == "OPEN"){
        this.openJobs.push(job);
      } else {
        this.closedJobs.push(job);
      }
    }
   }

  ngOnInit() {
  }

}
