import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModalModule, NgbDatepickerModule, NgbTimepickerModule } from '@ng-bootstrap/ng-bootstrap';

import { CreateJobRoutingModule } from "./create-job-routing.module";

import { CreateJobComponent } from "./create-job.component";



@NgModule({
  imports: [
    CommonModule,
    CreateJobRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    NgbModalModule,
    NgbDatepickerModule,
    NgbTimepickerModule,
  ],
  exports: [],
  declarations: [
    CreateJobComponent,
  ],
  providers: [],
})
export class CreateJobModule { }
