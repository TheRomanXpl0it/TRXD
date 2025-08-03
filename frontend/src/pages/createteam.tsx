import { useContext } from 'react';
import SettingContext from '@/context/SettingsProvider';
"use client"
import {
  useState
} from "react"
import {
  toast
} from "sonner"
import {
  useForm
} from "react-hook-form"
import {
  zodResolver
} from "@hookform/resolvers/zod"
import {
  z
} from "zod"
import {
  Button
} from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import {
  Input
} from "@/components/ui/input"
import LocationSelector from "@/components/ui/location-input"
import {
  PasswordInput
} from "@/components/ui/password-input"
import {
  CloudUpload,
  Paperclip
} from "lucide-react"
import {
  FileInput,
  FileUploader,
  FileUploaderContent,
  FileUploaderItem
} from "@/components/ui/file-upload"
import { registerTeam, updateTeam } from "@/lib/backend-interaction"
import { Label } from '@radix-ui/react-dropdown-menu';
import { register } from 'module';

const formSchema = z.object({
  TeamName: z.string().min(1),
  TeamDescription: z.string().optional(),
  TeamCountry: z.string().optional(),
  TeamPassword: z.string().min(8),
  TeamConfirmPassword: z.string().min(8),
  TeamProfilePicture: z.string().optional()
});

import AuthContext from '@/context/AuthProvider';

function CreateTeamForm() {
  const authContext = useContext(AuthContext);

  if (!authContext || !authContext.auth) {
    // You can render null, a spinner, or an error message
    return <div>Error: Auth context not available.</div>;
  }
  
  const [countryName, setCountryName] = useState < string > ('')
  const [files, setFiles] = useState < File[] | null > (null);
  // Removed usage of auth?.team as 'team' does not exist on AuthContextType
  const auth = authContext.auth;
  const team = auth.team || null;


  const dropZoneConfig = {
    maxFiles: 5,
    maxSize: 1024 * 1024 * 4,
    multiple: true,
  };
  const form = useForm < z.infer < typeof formSchema >> ({
    resolver: zodResolver(formSchema),

  })

  async function onSubmit(values: z.infer < typeof formSchema > ) {
    let teamId, teamName;
    try {
      const result = await registerTeam(values.TeamName, values.TeamPassword);
      switch (result.status) {
        case 200:
          teamId = result.data.id;
          teamName = result.data.name;
          break;
        case 400:
          toast.error("Password too short");
          form.setError("TeamPassword", {
            type: "manual",
            message: "Password must be at least 8 characters long"
          });
          return;
        case 409:
          const error = result.data.error;
          toast.error(error);
          form.setError("TeamName", {
            type: "manual",
            message: error,
          });
          return;
        case 500:
          toast.error("Team registration failed: Internal Server Error");
          return;
        default:
          toast.error("Team registration failed: Unknown error");
          return;
      }
    } catch (error) {
      console.error("Form submission error", error);
      toast.error("Failed to register the team.");
    }
    try {
      if ( values.TeamCountry || values.TeamDescription || values.TeamProfilePicture) {
        const result = await updateTeam(values.TeamDescription, values.TeamCountry, values.TeamProfilePicture);
        switch (result.status) {
          case 200:
            break;
          case 400:
            toast.error("Invalid data provided");
            form.setError("root", {
              type: "manual",
              message: "Invalid data provided"
            });
            return;
          case 500:
            toast.error("Team update failed: Internal Server Error");
            return;
          default:
            toast.error("Team update failed: Unknown error");
            return;
      }
    }
    } catch (error) {
      console.error("Form submission error", error);
      toast.error("Failed to update the team.");
    }
    
      authContext.setAuth({
        username: auth.username,
        roles: auth.roles,
        team: {
          id: teamId,
          name: teamName,
        }
      });
    toast.success("Team created successfully!");
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 max-w-3xl mx-auto py-10">
        
        <FormField
          control={form.control}
          name="TeamName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Team name*</FormLabel>
              <FormControl>
                <Input 
                placeholder="Team"
                type="text"
                {...field} />
              </FormControl>
              <FormDescription>The publicly displayed name of your team.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="TeamDescription"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Team Description</FormLabel>
              <FormControl>
                <Input 
                placeholder="Description"
                type="text"
                {...field} />
              </FormControl>
              <FormDescription>A brief description of your team.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
           <FormField
              control={form.control}
              name="TeamCountry"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Select Country</FormLabel>
                  <FormControl>
                  <LocationSelector
                    onCountryChange={(country) => {
                      setCountryName(country?.name || '')
                      form.setValue(field.name, country?.name || '')
                    }}
                  />
                  </FormControl>
                  <FormDescription>Select the team country</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
        
        <FormField
          control={form.control}
          name="TeamPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password*</FormLabel>
              <FormControl>
                <PasswordInput placeholder="Placeholder" {...field} />
              </FormControl>
              <FormDescription>Enter your password.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        
        <FormField
          control={form.control}
          name="TeamConfirmPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm Password*</FormLabel>
              <FormControl>
                <PasswordInput placeholder="Placeholder" {...field} />
              </FormControl>
              <FormDescription>Enter your password again.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        
        
            <FormField
              control={form.control}
              name="TeamProfilePicture"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Upload a team logo picture</FormLabel>
                  <FormControl>
                    <FileUploader
                      value={files}
                      onValueChange={setFiles}
                      dropzoneOptions={dropZoneConfig}
                      className="relative bg-background rounded-lg p-2"
                    >
                      <FileInput
                        id="fileInput"
                        className="outline-dashed outline-1 outline-slate-500"
                      >
                        <div className="flex items-center justify-center flex-col p-8 w-full ">
                          <CloudUpload className='text-gray-500 w-10 h-10' />
                          <p className="mb-1 text-sm text-gray-500 dark:text-gray-400">
                            <span className="font-semibold">Click to upload</span>
                            &nbsp; or drag and drop
                          </p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">
                            SVG, PNG, JPG or GIF
                          </p>
                        </div>
                      </FileInput>
                      <FileUploaderContent>
                        {files &&
                          files.length > 0 &&
                          files.map((file, i) => (
                            <FileUploaderItem key={i} index={i}>
                              <Paperclip className="h-4 w-4 stroke-current" />
                              <span>{file.name}</span>
                            </FileUploaderItem>
                          ))}
                      </FileUploaderContent>
                    </FileUploader>
                  </FormControl>
                  <FormDescription>Select a file to upload.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Label className="text-sm text-gray-500">
              All fields marked with * are required.
            </Label>
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}


export function CreateTeam() {
  const { settings } = useContext(SettingContext);
  const showQuotes = settings.General?.find((setting) => setting.title === 'Show Quotes')?.value;
  return (
    <>
       <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        Create a new Team
      </h2>
      { showQuotes && (
                <blockquote className="mt-6 border-l-2 pl-6 italic">
                    "Find a group of people who challenge and inspire you, spend a lot of time with them, and it will change your life."
                </blockquote>
            )}
      <CreateTeamForm />
    </>
  );
}
