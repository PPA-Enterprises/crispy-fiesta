import { Component, OnInit, Input, ChangeDetectorRef } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm, FormBuilder } from '@angular/forms';
import { JobService } from '../shared/services/job.service'
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { Job } from '../shared/models/job.model';
import { Color } from '../shared/models/color.model';
import { startOfDay, endOfDay, subDays, addDays, endOfMonth, isSameDay, isSameMonth, addHours } from 'date-fns';
import { NgbDateStruct, NgbTimeStruct, NgbCalendar, NgbTimepicker } from '@ng-bootstrap/ng-bootstrap';
import {
  getSeconds,
  getMinutes,
  getHours,
  getDate,
  getMonth,
  getYear,
  setSeconds,
  setMinutes,
  setHours,
  setDate,
  setMonth,
  setYear
} from 'date-fns';

export class JobForm {
  public fname: string;
  public lname: string;
  public phone: string;
  public carInfo: string;
  public apptInfo: string;
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
  @Input() placeholder: string;
  job: Job;
  model = new JobForm();
  submitted = false;
  startDateModel: NgbDateStruct;
  endDateModel: NgbDateStruct;
  startTimeModel = {hour: 12, minute: 0};
  endTimeModel = {hour: 12, minute: 0};
  startDate: Date;
  endDate: Date;
  form1: FormGroup;
  form2: FormGroup

  constructor(private jobService: JobService, private cdr: ChangeDetectorRef, private router: Router, private calendar: NgbCalendar, private fb: FormBuilder) { 
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
    // Use this method to set any other date format you want 
    this.startDate = new Date(date.year, date.month -1, date.day);
    console.log(this.startDate);
  }

  public updateEnd(date: NgbDateStruct) {
    // Use this method to set any other date format you want 
    this.endDate = new Date(date.year, date.month -1, date.day);
    console.log(this.endDate);
  }

  onSubmit(form) {
    console.log(this.startTimeModel);
    console.log(this.endTimeModel);
    console.log(this.startDate)
    
    
    this.submitted = true

    // this.job = {
    //   client_name: this.model.fname + " " + this.model.lname,
    //   client_phone: this.model.phone,
    //   car_info: this.model.carInfo,
    //   appointment_info: this.model.apptInfo,
    //   notes: this.model.notes,
    //   tag: "NEW",
    //   start: new Date(this.model.start.year, this.model.start.month - 1, this.model.start.day),
    //   end: new Date(this.model.end.year, this.model.end.month - 1, this.model.end.day),
    //   title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
    //   color: {
    //     primary: this.model.primary,
    //     secondary: this.model.secondary,
    //   }
    // }   

    // this.jobService.createJob(this.job).subscribe((job: Job) => {
    //     this.router.navigate(['/jobs']);
    // })
  }
}
