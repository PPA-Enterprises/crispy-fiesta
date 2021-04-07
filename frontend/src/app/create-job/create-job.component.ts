import { Component, OnInit, Input, ChangeDetectorRef } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm, FormBuilder } from '@angular/forms';
import { JobService } from '../shared/services/job.service'
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { Job } from '../shared/models/job.model';
import { Color } from '../shared/models/color.model';

export class JobForm {
  public fname: string;
  public lname: string;
  public phone: string;
  public carInfo: string;
  public notes: string;
  public start: any;
  public end: any;
  public primary: any;
  public secondary: any;
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
  startDate: any;
  endDate: any;

  constructor(private jobService: JobService, private router: Router) {}

  ngOnInit(): void {
  }


  onSubmit(form) {
    this.submitted = true

    this.job = {
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      notes: this.model.notes,
      tag: "NEW",
      start: this.startDate,
      end: this.endDate,
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      primary_color: this.model.primary,
      secondary_color: this.model.secondary,
    }   

    this.jobService.createJob(this.job).subscribe((job: Job) => {
        this.router.navigate(['/jobs']);
    })
  }
}
