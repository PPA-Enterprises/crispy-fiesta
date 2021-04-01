import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm, FormBuilder } from '@angular/forms';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { JobService } from '../shared/services/job.service'
import { Job } from '../shared/models/job.model';
import { startOfDay, endOfDay, subDays, addDays, endOfMonth, isSameDay, isSameMonth, addHours } from 'date-fns';
import { NgbDateStruct, NgbTimeStruct, NgbCalendar, NgbTimepicker } from '@ng-bootstrap/ng-bootstrap';

export class JobForm {
  public fname: string;
  public lname: string;
  public phone: string;
  public carInfo: string;
  public notes: string;
  public primary: any;
  public secondary: any;
}

@Component({
  selector: 'app-page',
  templateUrl: './edit-job.component.html',
  styleUrls: ['./edit-job.component.scss']
})

export class EditJobComponent {
  public job: Job;
  id: string;
  private sub: any;
  model = new JobForm();
  submitted = false;
  name: string[];
  notClosed = true;
  startTime = { hour: 12, minute: 30 };
  endTime = { hour: 12, minute: 30 };
  startDate = { year: 1789, month: 7, day: 14 };
  endDate = { year: 1789, month: 7, day: 14 };
  form1: FormGroup;
  form2: FormGroup;

  constructor(private jobService: JobService, private route: ActivatedRoute, private router: Router, private fb: FormBuilder, private calendar: NgbCalendar) {
    this.form1 = this.fb.group({
      'time' : [this.startTime, Validators.required],
    })
    this.form2 = this.fb.group({
      'time' : [this.endTime, Validators.required],
    })
    
    this.sub = this.route.params.subscribe(params => {
      this.id = params['id'];
    });

    this.jobService.getJobById(this.id).subscribe((incomingJob: Job) => {
      console.log(incomingJob);
      
      this.job = incomingJob;

      
      if(this.job.tag == "CLOSED") {
        this.notClosed = false;
      }
      this.name = this.job.client_name.split(" ");
      this.model = {
        fname: this.name[0],
        lname: this.name[1],
        phone: this.job.client_phone,
        carInfo: this.job.car_info,
        notes: this.job.notes,
        primary: this.job.color.primary,
        secondary: this.job.color.secondary
      }

      
      this.startTime.hour = this.job.start.getHours();
      this.startTime.minute = this.job.start.getMinutes();
      this.endTime.hour = this.job.end.getHours();
      this.endTime.minute = this.job.end.getMinutes();
      
      this.startDate.year = this.job.start.getFullYear();
      this.startDate.month = this.job.start.getMonth()+1;
      this.startDate.day = this.job.start.getDate();

      this.endDate.year = this.job.end.getFullYear();
      this.endDate.month = this.job.end.getMonth()+1;
      this.endDate.day = this.job.end.getDate();

      //this.calendar.startDate(this.startDate)
      
    })

    console.log(this.job);
  }

  ngOnInit() {
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }

  onSubmit(form) {
    this.job = {
      _id: this.job._id,
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      notes: this.model.notes,
      tag: "OPEN",
      start: new Date(this.startDate.year, this.startDate.month - 1, this.startDate.day, this.startTime.hour, this.startTime.minute),
      end: new Date(this.endDate.year, this.endDate.month - 1, this.endDate.day, this.endTime.hour, this.endTime.minute),
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      color: {
        primary: this.model.primary,
        secondary: this.model.secondary,
      }
    }

    this.jobService.editJobById(this.id, this.job).subscribe((job: Job) => {
      this.job = job
      this.router.navigate(['/jobs']);
    })
  }

  sendToTintWork() {
    console.log("Before:" + this.job);
    this.job = {
      _id: this.job._id,
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      notes: this.model.notes,
      tag: "CLOSED",
      start: subDays(startOfDay(new Date()), 1),
      end: addDays(new Date(), 1),
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      color: {
        primary: '#ad2121',
        secondary: '#FAE3E3',
      }
    }

    this.jobService.editJobById(this.id, this.job).subscribe((job: Job) => {
      this.job = job
      this.router.navigate(['/jobs']);
    })

  }

  delete() {
    this.jobService.deleteJobById(this.id).subscribe((deleted: any) => {
      if (deleted) this.router.navigate(['/jobs']);
    })
  }

  public updateStart(date: NgbDateStruct) {
    this.startDate = date;
  }

  public updateEnd(date: NgbDateStruct) {
    this.endDate = date;
  }


}
