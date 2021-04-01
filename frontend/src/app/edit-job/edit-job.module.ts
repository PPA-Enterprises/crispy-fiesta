import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModalModule, NgbDatepickerModule, NgbTimepickerModule } from '@ng-bootstrap/ng-bootstrap';

import { EditJobRoutingModule } from "./edit-job-routing.module";

import { EditJobComponent } from "./edit-job.component";


@NgModule({
  imports: [
    CommonModule,
    EditJobRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    NgbModalModule,
    NgbDatepickerModule,
    NgbTimepickerModule,
  ],
  exports: [],
  declarations: [
    EditJobComponent,
  ],
  providers: [],
})
export class EditJobModule { }
