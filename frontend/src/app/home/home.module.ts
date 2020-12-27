import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HomeRoutingModule } from "./home-routing.module";
import { NgxDatatableModule } from "@swimlane/ngx-datatable";
import { HomeComponent } from "./home.component";


@NgModule({
  imports: [
    CommonModule,
    HomeRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    NgxDatatableModule
  ],
  exports: [],
  declarations: [
    HomeComponent,
  ],
  providers: [],
})
export class HomeModule { }
