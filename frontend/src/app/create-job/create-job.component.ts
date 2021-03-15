import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { JobService } from '../shared/services/job.service'
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { Job } from '../shared/models/job.model';
import { Color } from '../shared/models/color.model';
import { startOfDay, endOfDay, subDays, addDays, endOfMonth, isSameDay, isSameMonth, addHours } from 'date-fns';


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
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      appointment_info: this.model.apptInfo,
      notes: this.model.notes,
      tag: "NEW",
      start: subDays(startOfDay(new Date()), 1),
      end: addDays(new Date(), 1),
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      color: {
        primary: '#ad2121',
        secondary: '#FAE3E3',
      }
    }

    this.jobService.createJob(this.job).subscribe((job: Job) => {
        this.router.navigate(['/jobs']);
    })
  }
}
