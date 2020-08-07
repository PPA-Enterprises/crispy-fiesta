import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { JobService } from '../shared/services/job.service'

export class JobForm {
  public fname: string;
  public lname: string;
  public phone: string;
  public carInfo: string;
  public apptInfo: string;
  public notes: string;
}

export class Job {
  public client_name: string
  public client_phone: string;
  public car_info: string;
  public appointment_info: string;
  public notes: string;
}

@Component({
  selector: 'app-page',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})

export class HomeComponent implements OnInit {
  job: Job;
  model = new JobForm();
  submitted = false;

  constructor(private jobService: JobService) {
    // this.model = {
    //   fname: 'Mark',
    //   lname: 'Otto',
    //   phone: ''
    // }
  }

  ngOnInit() {

  }

  onSubmit(form) {

    this.job = {
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      appointment_info: this.model.apptInfo,
      notes: this.model.notes
    }
    this.jobService.createJob(this.job).subscribe(data => { console.log(data)});
  }


}
