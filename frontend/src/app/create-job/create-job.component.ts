import { Component, OnInit, Input, ChangeDetectorRef } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm, FormBuilder } from '@angular/forms';
import { JobService } from '../shared/services/job.service'
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { Job } from '../shared/models/job.model';
import { Color } from '../shared/models/color.model';
import { NgbDateStruct, NgbTimeStruct, NgbCalendar, NgbTimepicker } from '@ng-bootstrap/ng-bootstrap';

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
  startDateModel: NgbDateStruct;
  endDateModel: NgbDateStruct;
  startTimeModel = {hour: 12, minute: 0};
  endTimeModel = {hour: 12, minute: 0};
  startDate: NgbDateStruct;
  endDate: NgbDateStruct;
  form1: FormGroup;
  form2: FormGroup

  constructor(private jobService: JobService, private router: Router, private calendar: NgbCalendar, private fb: FormBuilder) { 
    this.form1 = this.fb.group({
      'time' : [this.startTimeModel, Validators.required],
    })
    this.form2 = this.fb.group({
      'time' : [this.endTimeModel, Validators.required],
    })
  }

  ngOnInit(): void {
  }

  public updateStart(date: NgbDateStruct) {
    this.startDate = date;
  }

  public updateEnd(date: NgbDateStruct) {
    this.endDate = date;
  }

  onSubmit(form) {
    this.submitted = true

    this.job = {
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      notes: this.model.notes,
      tag: "NEW",
      start: new Date(this.startDate.year, this.startDate.month - 1, this.startDate.day, this.startTimeModel.hour, this.startTimeModel.minute),
      end: new Date(this.endDate.year, this.endDate.month - 1, this.endDate.day, this.endTimeModel.hour, this.endTimeModel.minute),
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      color: {
        primary: this.model.primary,
        secondary: this.model.secondary,
      }
    }   

    this.jobService.createJob(this.job).subscribe((job: Job) => {
        this.router.navigate(['/jobs']);
    })
  }
}
