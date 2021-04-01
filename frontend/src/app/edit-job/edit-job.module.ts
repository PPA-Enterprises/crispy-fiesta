import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { OwlDateTimeModule, OwlNativeDateTimeModule } from '@danielmoncada/angular-datetime-picker';

import { EditJobRoutingModule } from "./edit-job-routing.module";

import { EditJobComponent } from "./edit-job.component";


@NgModule({
  imports: [
    CommonModule,
    EditJobRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    OwlDateTimeModule,
    OwlNativeDateTimeModule,
  ],
  exports: [],
  declarations: [
    EditJobComponent,
  ],
  providers: [],
})
export class EditJobModule { }
