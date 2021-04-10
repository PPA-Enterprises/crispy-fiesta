import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { JobService } from '../shared/services/job.service'
import { Job } from '../shared/models/job.model'


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
  constructor(private jobService: JobService) {}

  ngOnInit() {
    this.jobService.getAllJobs().subscribe((jobs: Job[]) => {
      for(let job of jobs) {
        console.log(job)
        if(job.tag == "NEW") {
          this.newJobs.push(job);
        } else if(job.tag == "OPEN"){
          this.openJobs.push(job);
        } else {
          this.closedJobs.push(job);
        }
      }
    });
  }

}
