import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { JobService } from '../shared/services/job.service'
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

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

export class JobForm {
  public fname: string;
  public lname: string;
  public phone: string;
  public carInfo: string;
  public apptInfo: string;
  public notes: string;
}

@Component({
  selector: 'app-page',
  templateUrl: './create-job.component.html',
  styleUrls: ['./create-job.component.scss']
})
export class CreateJobComponent implements OnInit {
  job: Job;
  model = new JobForm();
  submitted = false;
  
  constructor(private jobService: JobService, private router: Router) { }

  ngOnInit(): void {
  }

  onSubmit(form) {
    this.job = {
      id: Math.floor(Math.random() * Math.floor(250)),
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      appointment_info: this.model.apptInfo,
      notes: this.model.notes,
      tag: "NEW",
      date: new Date().toLocaleString(),
    }
    
    if (this.jobService.createJob(this.job).tag == "NEW") {
      this.router.navigate(['/jobs']);
    } else {
      console.log("ERROR CREATING JOB!")
    }
    
  }

}
