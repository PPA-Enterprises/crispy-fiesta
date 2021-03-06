import { Routes, RouterModule } from '@angular/router';

//Route for content layout with sidebar, navbar and footer.

export const Full_ROUTES: Routes = [
  {
    path: 'home',
    loadChildren: () => import('../../home/home.module').then(m => m.HomeModule)
  },
  {
    path: 'jobs',
    loadChildren: () => import('../../jobs/jobs.module').then(m => m.JobsModule)
  },
  {
    path: 'clients',
    loadChildren: () => import('../../clients/clients.module').then(m => m.ClientsModule)
  },
  {
    path: 'client/:id',
    loadChildren: () => import('../../client/client.module').then(m => m.ClientModule)
  },
  {
    path: 'create-job',
    loadChildren: () => import('../../create-job/create-job.module').then(m => m.CreateJobModule)
  },
  {
    path: 'edit-job/:id',
    loadChildren: () => import('../../edit-job/edit-job.module').then(m => m.EditJobModule)
  },
  {
    path: 'calendar',
    loadChildren: () => import('../../cal/cal.module').then(m => m.CalModule)
  },
];
