import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { EditJobComponent } from './edit-job.component';


const routes: Routes = [
  {
    path: '',
    component: EditJobComponent,
    data: {
      title: 'Edit Job'
    }
    // children: [
    //   {
    //     path: 'page',
    //     component: PageComponent,
    //     data: {
    //       title: 'Page'
    //     }
    //   }
    // ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class EditJobRoutingModule { }
