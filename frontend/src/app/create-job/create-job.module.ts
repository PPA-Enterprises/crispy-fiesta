import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { OwlDateTimeModule, OwlNativeDateTimeModule } from '@danielmoncada/angular-datetime-picker';

import { CreateJobRoutingModule } from "./create-job-routing.module";

import { CreateJobComponent } from "./create-job.component";



@NgModule({
  imports: [
    CommonModule,
    CreateJobRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    OwlDateTimeModule,
    OwlNativeDateTimeModule,
  ],
  exports: [],
  declarations: [
    CreateJobComponent,
  ],
  providers: [],
})
export class CreateJobModule { }
